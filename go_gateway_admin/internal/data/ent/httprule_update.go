// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/baker-yuan/go-gateway/go-gateway-admin/internal/data/ent/httprule"
	"github.com/baker-yuan/go-gateway/go-gateway-admin/internal/data/ent/predicate"
)

// HttpRuleUpdate is the builder for updating HttpRule entities.
type HttpRuleUpdate struct {
	config
	hooks    []Hook
	mutation *HttpRuleMutation
}

// Where appends a list predicates to the HttpRuleUpdate builder.
func (hru *HttpRuleUpdate) Where(ps ...predicate.HttpRule) *HttpRuleUpdate {
	hru.mutation.Where(ps...)
	return hru
}

// SetGwURL sets the "gw_url" field.
func (hru *HttpRuleUpdate) SetGwURL(s string) *HttpRuleUpdate {
	hru.mutation.SetGwURL(s)
	return hru
}

// SetNillableGwURL sets the "gw_url" field if the given value is not nil.
func (hru *HttpRuleUpdate) SetNillableGwURL(s *string) *HttpRuleUpdate {
	if s != nil {
		hru.SetGwURL(*s)
	}
	return hru
}

// SetHTTPType sets the "http_type" field.
func (hru *HttpRuleUpdate) SetHTTPType(s string) *HttpRuleUpdate {
	hru.mutation.SetHTTPType(s)
	return hru
}

// SetNillableHTTPType sets the "http_type" field if the given value is not nil.
func (hru *HttpRuleUpdate) SetNillableHTTPType(s *string) *HttpRuleUpdate {
	if s != nil {
		hru.SetHTTPType(*s)
	}
	return hru
}

// SetStatus sets the "status" field.
func (hru *HttpRuleUpdate) SetStatus(u uint8) *HttpRuleUpdate {
	hru.mutation.ResetStatus()
	hru.mutation.SetStatus(u)
	return hru
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (hru *HttpRuleUpdate) SetNillableStatus(u *uint8) *HttpRuleUpdate {
	if u != nil {
		hru.SetStatus(*u)
	}
	return hru
}

// AddStatus adds u to the "status" field.
func (hru *HttpRuleUpdate) AddStatus(u int8) *HttpRuleUpdate {
	hru.mutation.AddStatus(u)
	return hru
}

// SetApplication sets the "application" field.
func (hru *HttpRuleUpdate) SetApplication(s string) *HttpRuleUpdate {
	hru.mutation.SetApplication(s)
	return hru
}

// SetInterfaceType sets the "interface_type" field.
func (hru *HttpRuleUpdate) SetInterfaceType(u uint8) *HttpRuleUpdate {
	hru.mutation.ResetInterfaceType()
	hru.mutation.SetInterfaceType(u)
	return hru
}

// SetNillableInterfaceType sets the "interface_type" field if the given value is not nil.
func (hru *HttpRuleUpdate) SetNillableInterfaceType(u *uint8) *HttpRuleUpdate {
	if u != nil {
		hru.SetInterfaceType(*u)
	}
	return hru
}

// AddInterfaceType adds u to the "interface_type" field.
func (hru *HttpRuleUpdate) AddInterfaceType(u int8) *HttpRuleUpdate {
	hru.mutation.AddInterfaceType(u)
	return hru
}

// SetInterfaceURL sets the "interface_url" field.
func (hru *HttpRuleUpdate) SetInterfaceURL(s string) *HttpRuleUpdate {
	hru.mutation.SetInterfaceURL(s)
	return hru
}

// SetNillableInterfaceURL sets the "interface_url" field if the given value is not nil.
func (hru *HttpRuleUpdate) SetNillableInterfaceURL(s *string) *HttpRuleUpdate {
	if s != nil {
		hru.SetInterfaceURL(*s)
	}
	return hru
}

// SetConfig sets the "config" field.
func (hru *HttpRuleUpdate) SetConfig(s string) *HttpRuleUpdate {
	hru.mutation.SetConfig(s)
	return hru
}

// SetNillableConfig sets the "config" field if the given value is not nil.
func (hru *HttpRuleUpdate) SetNillableConfig(s *string) *HttpRuleUpdate {
	if s != nil {
		hru.SetConfig(*s)
	}
	return hru
}

// SetCreateTime sets the "create_time" field.
func (hru *HttpRuleUpdate) SetCreateTime(u uint32) *HttpRuleUpdate {
	hru.mutation.ResetCreateTime()
	hru.mutation.SetCreateTime(u)
	return hru
}

// SetNillableCreateTime sets the "create_time" field if the given value is not nil.
func (hru *HttpRuleUpdate) SetNillableCreateTime(u *uint32) *HttpRuleUpdate {
	if u != nil {
		hru.SetCreateTime(*u)
	}
	return hru
}

// AddCreateTime adds u to the "create_time" field.
func (hru *HttpRuleUpdate) AddCreateTime(u int32) *HttpRuleUpdate {
	hru.mutation.AddCreateTime(u)
	return hru
}

// SetUpdateTime sets the "update_time" field.
func (hru *HttpRuleUpdate) SetUpdateTime(u uint32) *HttpRuleUpdate {
	hru.mutation.ResetUpdateTime()
	hru.mutation.SetUpdateTime(u)
	return hru
}

// SetNillableUpdateTime sets the "update_time" field if the given value is not nil.
func (hru *HttpRuleUpdate) SetNillableUpdateTime(u *uint32) *HttpRuleUpdate {
	if u != nil {
		hru.SetUpdateTime(*u)
	}
	return hru
}

// AddUpdateTime adds u to the "update_time" field.
func (hru *HttpRuleUpdate) AddUpdateTime(u int32) *HttpRuleUpdate {
	hru.mutation.AddUpdateTime(u)
	return hru
}

// Mutation returns the HttpRuleMutation object of the builder.
func (hru *HttpRuleUpdate) Mutation() *HttpRuleMutation {
	return hru.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (hru *HttpRuleUpdate) Save(ctx context.Context) (int, error) {
	return withHooks[int, HttpRuleMutation](ctx, hru.sqlSave, hru.mutation, hru.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (hru *HttpRuleUpdate) SaveX(ctx context.Context) int {
	affected, err := hru.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (hru *HttpRuleUpdate) Exec(ctx context.Context) error {
	_, err := hru.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (hru *HttpRuleUpdate) ExecX(ctx context.Context) {
	if err := hru.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (hru *HttpRuleUpdate) check() error {
	if v, ok := hru.mutation.GwURL(); ok {
		if err := httprule.GwURLValidator(v); err != nil {
			return &ValidationError{Name: "gw_url", err: fmt.Errorf(`ent: validator failed for field "HttpRule.gw_url": %w`, err)}
		}
	}
	if v, ok := hru.mutation.HTTPType(); ok {
		if err := httprule.HTTPTypeValidator(v); err != nil {
			return &ValidationError{Name: "http_type", err: fmt.Errorf(`ent: validator failed for field "HttpRule.http_type": %w`, err)}
		}
	}
	if v, ok := hru.mutation.Application(); ok {
		if err := httprule.ApplicationValidator(v); err != nil {
			return &ValidationError{Name: "application", err: fmt.Errorf(`ent: validator failed for field "HttpRule.application": %w`, err)}
		}
	}
	if v, ok := hru.mutation.InterfaceURL(); ok {
		if err := httprule.InterfaceURLValidator(v); err != nil {
			return &ValidationError{Name: "interface_url", err: fmt.Errorf(`ent: validator failed for field "HttpRule.interface_url": %w`, err)}
		}
	}
	if v, ok := hru.mutation.Config(); ok {
		if err := httprule.ConfigValidator(v); err != nil {
			return &ValidationError{Name: "config", err: fmt.Errorf(`ent: validator failed for field "HttpRule.config": %w`, err)}
		}
	}
	return nil
}

func (hru *HttpRuleUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := hru.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(httprule.Table, httprule.Columns, sqlgraph.NewFieldSpec(httprule.FieldID, field.TypeUint32))
	if ps := hru.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := hru.mutation.GwURL(); ok {
		_spec.SetField(httprule.FieldGwURL, field.TypeString, value)
	}
	if value, ok := hru.mutation.HTTPType(); ok {
		_spec.SetField(httprule.FieldHTTPType, field.TypeString, value)
	}
	if value, ok := hru.mutation.Status(); ok {
		_spec.SetField(httprule.FieldStatus, field.TypeUint8, value)
	}
	if value, ok := hru.mutation.AddedStatus(); ok {
		_spec.AddField(httprule.FieldStatus, field.TypeUint8, value)
	}
	if value, ok := hru.mutation.Application(); ok {
		_spec.SetField(httprule.FieldApplication, field.TypeString, value)
	}
	if value, ok := hru.mutation.InterfaceType(); ok {
		_spec.SetField(httprule.FieldInterfaceType, field.TypeUint8, value)
	}
	if value, ok := hru.mutation.AddedInterfaceType(); ok {
		_spec.AddField(httprule.FieldInterfaceType, field.TypeUint8, value)
	}
	if value, ok := hru.mutation.InterfaceURL(); ok {
		_spec.SetField(httprule.FieldInterfaceURL, field.TypeString, value)
	}
	if value, ok := hru.mutation.Config(); ok {
		_spec.SetField(httprule.FieldConfig, field.TypeString, value)
	}
	if value, ok := hru.mutation.CreateTime(); ok {
		_spec.SetField(httprule.FieldCreateTime, field.TypeUint32, value)
	}
	if value, ok := hru.mutation.AddedCreateTime(); ok {
		_spec.AddField(httprule.FieldCreateTime, field.TypeUint32, value)
	}
	if value, ok := hru.mutation.UpdateTime(); ok {
		_spec.SetField(httprule.FieldUpdateTime, field.TypeUint32, value)
	}
	if value, ok := hru.mutation.AddedUpdateTime(); ok {
		_spec.AddField(httprule.FieldUpdateTime, field.TypeUint32, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, hru.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{httprule.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	hru.mutation.done = true
	return n, nil
}

// HttpRuleUpdateOne is the builder for updating a single HttpRule entity.
type HttpRuleUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *HttpRuleMutation
}

// SetGwURL sets the "gw_url" field.
func (hruo *HttpRuleUpdateOne) SetGwURL(s string) *HttpRuleUpdateOne {
	hruo.mutation.SetGwURL(s)
	return hruo
}

// SetNillableGwURL sets the "gw_url" field if the given value is not nil.
func (hruo *HttpRuleUpdateOne) SetNillableGwURL(s *string) *HttpRuleUpdateOne {
	if s != nil {
		hruo.SetGwURL(*s)
	}
	return hruo
}

// SetHTTPType sets the "http_type" field.
func (hruo *HttpRuleUpdateOne) SetHTTPType(s string) *HttpRuleUpdateOne {
	hruo.mutation.SetHTTPType(s)
	return hruo
}

// SetNillableHTTPType sets the "http_type" field if the given value is not nil.
func (hruo *HttpRuleUpdateOne) SetNillableHTTPType(s *string) *HttpRuleUpdateOne {
	if s != nil {
		hruo.SetHTTPType(*s)
	}
	return hruo
}

// SetStatus sets the "status" field.
func (hruo *HttpRuleUpdateOne) SetStatus(u uint8) *HttpRuleUpdateOne {
	hruo.mutation.ResetStatus()
	hruo.mutation.SetStatus(u)
	return hruo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (hruo *HttpRuleUpdateOne) SetNillableStatus(u *uint8) *HttpRuleUpdateOne {
	if u != nil {
		hruo.SetStatus(*u)
	}
	return hruo
}

// AddStatus adds u to the "status" field.
func (hruo *HttpRuleUpdateOne) AddStatus(u int8) *HttpRuleUpdateOne {
	hruo.mutation.AddStatus(u)
	return hruo
}

// SetApplication sets the "application" field.
func (hruo *HttpRuleUpdateOne) SetApplication(s string) *HttpRuleUpdateOne {
	hruo.mutation.SetApplication(s)
	return hruo
}

// SetInterfaceType sets the "interface_type" field.
func (hruo *HttpRuleUpdateOne) SetInterfaceType(u uint8) *HttpRuleUpdateOne {
	hruo.mutation.ResetInterfaceType()
	hruo.mutation.SetInterfaceType(u)
	return hruo
}

// SetNillableInterfaceType sets the "interface_type" field if the given value is not nil.
func (hruo *HttpRuleUpdateOne) SetNillableInterfaceType(u *uint8) *HttpRuleUpdateOne {
	if u != nil {
		hruo.SetInterfaceType(*u)
	}
	return hruo
}

// AddInterfaceType adds u to the "interface_type" field.
func (hruo *HttpRuleUpdateOne) AddInterfaceType(u int8) *HttpRuleUpdateOne {
	hruo.mutation.AddInterfaceType(u)
	return hruo
}

// SetInterfaceURL sets the "interface_url" field.
func (hruo *HttpRuleUpdateOne) SetInterfaceURL(s string) *HttpRuleUpdateOne {
	hruo.mutation.SetInterfaceURL(s)
	return hruo
}

// SetNillableInterfaceURL sets the "interface_url" field if the given value is not nil.
func (hruo *HttpRuleUpdateOne) SetNillableInterfaceURL(s *string) *HttpRuleUpdateOne {
	if s != nil {
		hruo.SetInterfaceURL(*s)
	}
	return hruo
}

// SetConfig sets the "config" field.
func (hruo *HttpRuleUpdateOne) SetConfig(s string) *HttpRuleUpdateOne {
	hruo.mutation.SetConfig(s)
	return hruo
}

// SetNillableConfig sets the "config" field if the given value is not nil.
func (hruo *HttpRuleUpdateOne) SetNillableConfig(s *string) *HttpRuleUpdateOne {
	if s != nil {
		hruo.SetConfig(*s)
	}
	return hruo
}

// SetCreateTime sets the "create_time" field.
func (hruo *HttpRuleUpdateOne) SetCreateTime(u uint32) *HttpRuleUpdateOne {
	hruo.mutation.ResetCreateTime()
	hruo.mutation.SetCreateTime(u)
	return hruo
}

// SetNillableCreateTime sets the "create_time" field if the given value is not nil.
func (hruo *HttpRuleUpdateOne) SetNillableCreateTime(u *uint32) *HttpRuleUpdateOne {
	if u != nil {
		hruo.SetCreateTime(*u)
	}
	return hruo
}

// AddCreateTime adds u to the "create_time" field.
func (hruo *HttpRuleUpdateOne) AddCreateTime(u int32) *HttpRuleUpdateOne {
	hruo.mutation.AddCreateTime(u)
	return hruo
}

// SetUpdateTime sets the "update_time" field.
func (hruo *HttpRuleUpdateOne) SetUpdateTime(u uint32) *HttpRuleUpdateOne {
	hruo.mutation.ResetUpdateTime()
	hruo.mutation.SetUpdateTime(u)
	return hruo
}

// SetNillableUpdateTime sets the "update_time" field if the given value is not nil.
func (hruo *HttpRuleUpdateOne) SetNillableUpdateTime(u *uint32) *HttpRuleUpdateOne {
	if u != nil {
		hruo.SetUpdateTime(*u)
	}
	return hruo
}

// AddUpdateTime adds u to the "update_time" field.
func (hruo *HttpRuleUpdateOne) AddUpdateTime(u int32) *HttpRuleUpdateOne {
	hruo.mutation.AddUpdateTime(u)
	return hruo
}

// Mutation returns the HttpRuleMutation object of the builder.
func (hruo *HttpRuleUpdateOne) Mutation() *HttpRuleMutation {
	return hruo.mutation
}

// Where appends a list predicates to the HttpRuleUpdate builder.
func (hruo *HttpRuleUpdateOne) Where(ps ...predicate.HttpRule) *HttpRuleUpdateOne {
	hruo.mutation.Where(ps...)
	return hruo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (hruo *HttpRuleUpdateOne) Select(field string, fields ...string) *HttpRuleUpdateOne {
	hruo.fields = append([]string{field}, fields...)
	return hruo
}

// Save executes the query and returns the updated HttpRule entity.
func (hruo *HttpRuleUpdateOne) Save(ctx context.Context) (*HttpRule, error) {
	return withHooks[*HttpRule, HttpRuleMutation](ctx, hruo.sqlSave, hruo.mutation, hruo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (hruo *HttpRuleUpdateOne) SaveX(ctx context.Context) *HttpRule {
	node, err := hruo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (hruo *HttpRuleUpdateOne) Exec(ctx context.Context) error {
	_, err := hruo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (hruo *HttpRuleUpdateOne) ExecX(ctx context.Context) {
	if err := hruo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (hruo *HttpRuleUpdateOne) check() error {
	if v, ok := hruo.mutation.GwURL(); ok {
		if err := httprule.GwURLValidator(v); err != nil {
			return &ValidationError{Name: "gw_url", err: fmt.Errorf(`ent: validator failed for field "HttpRule.gw_url": %w`, err)}
		}
	}
	if v, ok := hruo.mutation.HTTPType(); ok {
		if err := httprule.HTTPTypeValidator(v); err != nil {
			return &ValidationError{Name: "http_type", err: fmt.Errorf(`ent: validator failed for field "HttpRule.http_type": %w`, err)}
		}
	}
	if v, ok := hruo.mutation.Application(); ok {
		if err := httprule.ApplicationValidator(v); err != nil {
			return &ValidationError{Name: "application", err: fmt.Errorf(`ent: validator failed for field "HttpRule.application": %w`, err)}
		}
	}
	if v, ok := hruo.mutation.InterfaceURL(); ok {
		if err := httprule.InterfaceURLValidator(v); err != nil {
			return &ValidationError{Name: "interface_url", err: fmt.Errorf(`ent: validator failed for field "HttpRule.interface_url": %w`, err)}
		}
	}
	if v, ok := hruo.mutation.Config(); ok {
		if err := httprule.ConfigValidator(v); err != nil {
			return &ValidationError{Name: "config", err: fmt.Errorf(`ent: validator failed for field "HttpRule.config": %w`, err)}
		}
	}
	return nil
}

func (hruo *HttpRuleUpdateOne) sqlSave(ctx context.Context) (_node *HttpRule, err error) {
	if err := hruo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(httprule.Table, httprule.Columns, sqlgraph.NewFieldSpec(httprule.FieldID, field.TypeUint32))
	id, ok := hruo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "HttpRule.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := hruo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, httprule.FieldID)
		for _, f := range fields {
			if !httprule.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != httprule.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := hruo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := hruo.mutation.GwURL(); ok {
		_spec.SetField(httprule.FieldGwURL, field.TypeString, value)
	}
	if value, ok := hruo.mutation.HTTPType(); ok {
		_spec.SetField(httprule.FieldHTTPType, field.TypeString, value)
	}
	if value, ok := hruo.mutation.Status(); ok {
		_spec.SetField(httprule.FieldStatus, field.TypeUint8, value)
	}
	if value, ok := hruo.mutation.AddedStatus(); ok {
		_spec.AddField(httprule.FieldStatus, field.TypeUint8, value)
	}
	if value, ok := hruo.mutation.Application(); ok {
		_spec.SetField(httprule.FieldApplication, field.TypeString, value)
	}
	if value, ok := hruo.mutation.InterfaceType(); ok {
		_spec.SetField(httprule.FieldInterfaceType, field.TypeUint8, value)
	}
	if value, ok := hruo.mutation.AddedInterfaceType(); ok {
		_spec.AddField(httprule.FieldInterfaceType, field.TypeUint8, value)
	}
	if value, ok := hruo.mutation.InterfaceURL(); ok {
		_spec.SetField(httprule.FieldInterfaceURL, field.TypeString, value)
	}
	if value, ok := hruo.mutation.Config(); ok {
		_spec.SetField(httprule.FieldConfig, field.TypeString, value)
	}
	if value, ok := hruo.mutation.CreateTime(); ok {
		_spec.SetField(httprule.FieldCreateTime, field.TypeUint32, value)
	}
	if value, ok := hruo.mutation.AddedCreateTime(); ok {
		_spec.AddField(httprule.FieldCreateTime, field.TypeUint32, value)
	}
	if value, ok := hruo.mutation.UpdateTime(); ok {
		_spec.SetField(httprule.FieldUpdateTime, field.TypeUint32, value)
	}
	if value, ok := hruo.mutation.AddedUpdateTime(); ok {
		_spec.AddField(httprule.FieldUpdateTime, field.TypeUint32, value)
	}
	_node = &HttpRule{config: hruo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, hruo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{httprule.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	hruo.mutation.done = true
	return _node, nil
}