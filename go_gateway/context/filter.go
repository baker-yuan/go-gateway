package eocontext

import "github.com/baker-yuan/go-gateway/utils/config"

var (
	FilterSkillName = config.TypeNameOf((*IFilter)(nil))
)

// IFilter 过滤器
type IFilter interface {
	DoFilter(ctx EoContext, next IChain) (err error) // 过滤逻辑
	Destroy()                                        // 过滤器销毁逻辑
}

// IChain 拦击器链
type IChain interface {
	DoChain(ctx EoContext) error // 往下执行
	Destroy()                    // 销毁
}

// IChainPro 拦击器链升级版本
type IChainPro interface {
	Chain(ctx EoContext, append ...IFilter) error //
	Destroy()                                     // 销毁
}

type Filters []IFilter

func (fs Filters) DoChain(ctx EoContext) error {
	if len(fs) > 0 {
		// 把链上Filters第一拿出来，改变Filters，扔掉第一，执行
		f := fs[0]
		next := fs[1:]
		return f.DoFilter(ctx, next) // Filter执行DoFilter方法next.DoChain(ctx)，又会到这里来，重复这个逻辑，知道这个链条每个Filter都被执行了一遍
	}
	return nil
}

func (fs Filters) Destroy() {
	for _, f := range fs {
		f.Destroy()
	}
}

func DoChain(ctx EoContext, orgfilter Filters, append ...IFilter) error {
	fs := orgfilter
	fl := len(fs)
	al := len(append)
	if fl == 0 && al == 0 {
		return nil
	}
	if fl == 0 {
		return Filters(append).DoChain(ctx)
	}
	if al == 0 {
		return fs.DoChain(ctx)
	}

	tp := make(Filters, fl+al)
	copy(tp, fs)
	copy(tp[fl:], append)
	return tp.DoChain(ctx)
}

type _FilterChain struct {
	chain IChain
}

func (c *_FilterChain) DoFilter(ctx EoContext, next IChain) (err error) {
	return c.chain.DoChain(ctx)
}

func (c *_FilterChain) Destroy() {
	c.chain.Destroy()
}

func ToFilter(chain IChain) IFilter {
	return &_FilterChain{chain: chain}
}
