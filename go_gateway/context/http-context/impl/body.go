package http_context

import (
	"bytes"
	"io"
	"mime"
	"mime/multipart"
	"net/url"
	"strings"

	http_context "github.com/baker-yuan/go-gateway/context/http-context"
	"github.com/valyala/fasthttp"
)

const defaultMultipartMemory = 32 << 20 // 32 MB

var (
	_ http_context.IBodyDataWriter = (*BodyRequestHandler)(nil)
)

// HTTP协议中的内容类型（Content-Type）
const (
	// MultipartForm
	// 主要用于在表单中发送二进制数据。最常见的用例是上传文件。在此编码类型中，表单的每个字段被视为一部分（multipart），每个部分都包含有关该字段的信息，例如字段名，字段值，如果字段是文件，则还包含文件名和文件类型。这意味着，使用这种类型，你可以在同一请求中发送文本和数据。
	MultipartForm = "multipart/form-data"
	// FormData
	// 通常用于发送ASCII字符集的数据。在此编码类型中，表单的字段名和值用等号（=）连接，字段之间用&符号分隔。所有非字母数字字符都会被百分比编码。这种类型常用于提交简单的文本数据。
	FormData       = "application/x-www-form-urlencoded"
	TEXT           = "text/plain"
	JSON           = "application/json"
	JavaScript     = "application/javascript"
	AppLicationXML = "application/xml"
	TextXML        = "text/xml"
	Html           = "text/html"
)

// BodyRequestHandler body请求处理器
// go_gateway/context/http-context/context.go#IBodyDataReader
// go_gateway/context/http-context/context.go#IBodyDataWriter
type BodyRequestHandler struct {
	request  *fasthttp.Request
	formdata *multipart.Form
}

func NewBodyRequestHandler(request *fasthttp.Request) *BodyRequestHandler {
	return &BodyRequestHandler{request: request}
}

// --------------------------- go_gateway/context/http-context/context.go#IBodyDataReader

// ContentType 获取contentType
func (b *BodyRequestHandler) ContentType() string {
	return string(b.request.Header.ContentType())
}

// BodyForm 获取表单参数
func (b *BodyRequestHandler) BodyForm() (url.Values, error) {
	contentType, _, _ := mime.ParseMediaType(string(b.request.Header.ContentType()))
	switch contentType {
	case FormData:
		return url.ParseQuery(string(b.request.Body()))
	case MultipartForm:
		multipartForm, err := b.MultipartForm()
		if err != nil {
			return nil, err
		}
		return multipartForm.Value, nil
	default:
		return nil, ErrorNotForm
	}

}

func (b *BodyRequestHandler) Files() (map[string][]*multipart.FileHeader, error) {
	form, err := b.MultipartForm()
	if err != nil {
		return nil, err
	}
	return form.File, nil
}

// GetForm 获取表单参数
func (b *BodyRequestHandler) GetForm(key string) string {
	contentType, _, _ := mime.ParseMediaType(b.ContentType())
	switch contentType {
	case FormData:
		args := b.request.PostArgs()
		if args == nil {
			return ""
		}
		return string(args.Peek(key))
	case MultipartForm:
		form, err := b.MultipartForm()
		if err != nil {
			return ""
		}
		vs := form.Value[key]
		if len(vs) > 0 {
			return vs[0]
		}
		return ""
	}
	return ""
}

func (b *BodyRequestHandler) GetFile(key string) ([]*multipart.FileHeader, bool) {
	multipartForm, err := b.MultipartForm()
	if err != nil {
		return nil, false
	}
	fl, has := multipartForm.File[key]
	return fl, has
}

// RawBody 获取raw数据
func (b *BodyRequestHandler) RawBody() ([]byte, error) {
	return b.request.Body(), nil
}

// --------------------------- go_gateway/context/http-context/context.go#IBodyDataWriter

// SetForm 设置表单参数
func (b *BodyRequestHandler) SetForm(values url.Values) error {
	contentType, _, _ := mime.ParseMediaType(b.ContentType())
	if contentType != FormData && contentType != MultipartForm {
		return ErrorNotForm
	}
	switch contentType {
	case FormData:
		b.request.SetBodyString(values.Encode())
	case MultipartForm:
		multipartForm, err := b.MultipartForm()
		if err != nil {
			return err
		}
		multipartForm.Value = values
		return b.resetFile()
	}
	return ErrorNotForm
}

func (b *BodyRequestHandler) SetToForm(key, value string) error {
	contentType, _, _ := mime.ParseMediaType(string(b.request.Header.ContentType()))
	switch contentType {
	case FormData:
		b.request.PostArgs().Set(key, value)
		b.request.SetBodyRaw(b.request.PostArgs().QueryString())
		return nil
	case MultipartForm:
		multipartForm, err := b.MultipartForm()
		if err != nil {
			return err
		}
		multipartForm.Value[key] = []string{value}
		return b.resetFile()
	default:
		return ErrorNotForm
	}
}

// AddForm 新增表单参数
func (b *BodyRequestHandler) AddForm(key, value string) error {
	contentType, _, _ := mime.ParseMediaType(string(b.request.Header.ContentType()))
	switch contentType {
	case FormData:
		b.request.PostArgs().Add(key, value)
		b.request.SetBody(b.request.PostArgs().QueryString())
		return nil
	case MultipartForm:
		multipartForm, err := b.MultipartForm()
		if err != nil {
			return err
		}
		multipartForm.Value[key] = append(multipartForm.Value[key], value)
		return b.resetFile()
	default:
		return ErrorNotForm
	}
}

// AddFile 新增文件参数
func (b *BodyRequestHandler) AddFile(key string, file *multipart.FileHeader) error {
	contentType, _, _ := mime.ParseMediaType(b.ContentType())
	if contentType != FormData && contentType != MultipartForm {
		return ErrorNotMultipart
	}
	multipartForm, err := b.MultipartForm()
	if err != nil {
		return err
	}
	multipartForm.File[key] = append(multipartForm.File[key], file)
	return b.resetFile()
}

// SetRaw 设置raw数据
func (b *BodyRequestHandler) SetRaw(contentType string, body []byte) {
	b.request.SetBodyRaw(body)
	b.request.Header.SetContentType(contentType)
}

// ---------------------------

// SetFile 设置文件参数
func (b *BodyRequestHandler) SetFile(files map[string][]*multipart.FileHeader) error {
	multipartForm, err := b.MultipartForm()
	if err != nil {
		return err
	}
	multipartForm.File = files
	return b.resetFile()
}

func (b *BodyRequestHandler) MultipartForm() (*multipart.Form, error) {
	if b.formdata != nil {
		return b.formdata, nil
	}
	if !strings.Contains(b.ContentType(), MultipartForm) {
		return nil, ErrorNotMultipart
	}
	form, err := b.request.MultipartForm()
	if err != nil {
		return nil, err
	}

	b.formdata = &multipart.Form{
		Value: form.Value,
		File:  form.File,
	}
	b.resetFile()
	return form, nil
}

func (b *BodyRequestHandler) resetFile() error {
	multipartForm := b.formdata
	if multipartForm == nil {
		return nil
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	for name, fs := range multipartForm.File {
		for _, f := range fs {
			fio, err := f.Open()
			if err != nil {
				return err
			}

			part, err := writer.CreateFormFile(name, f.Filename)
			if err != nil {
				fio.Close()
				return err
			}

			data, err := io.ReadAll(fio)
			if err != nil {
				fio.Close()
				return err
			}
			_, err = part.Write(data)
			if err != nil {
				fio.Close()
				return err
			}
			fio.Close()
		}
	}

	for key, values := range multipartForm.Value {
		// temp := make(url.Values)
		// temp[key] = values
		// value := temp.Encode()
		for _, value := range values {
			err := writer.WriteField(key, value)
			if err != nil {
				return err
			}
		}

	}
	err := writer.Close()
	if err != nil {
		return err
	}
	b.request.Header.SetContentType(writer.FormDataContentType())
	b.request.SetBodyRaw(body.Bytes())
	return nil
}

func (b *BodyRequestHandler) reset(request *fasthttp.Request) {
	b.request = request
	b.formdata = nil
}
