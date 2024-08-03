// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/member-gentei/member-gentei/gentei/ent/guild"
	"github.com/member-gentei/member-gentei/gentei/ent/guildrole"
	"github.com/member-gentei/member-gentei/gentei/ent/predicate"
	"github.com/member-gentei/member-gentei/gentei/ent/usermembership"
	"github.com/member-gentei/member-gentei/gentei/ent/youtubetalent"
)

// GuildRoleQuery is the builder for querying GuildRole entities.
type GuildRoleQuery struct {
	config
	ctx                 *QueryContext
	order               []guildrole.OrderOption
	inters              []Interceptor
	predicates          []predicate.GuildRole
	withGuild           *GuildQuery
	withUserMemberships *UserMembershipQuery
	withTalent          *YouTubeTalentQuery
	withFKs             bool
	modifiers           []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the GuildRoleQuery builder.
func (grq *GuildRoleQuery) Where(ps ...predicate.GuildRole) *GuildRoleQuery {
	grq.predicates = append(grq.predicates, ps...)
	return grq
}

// Limit the number of records to be returned by this query.
func (grq *GuildRoleQuery) Limit(limit int) *GuildRoleQuery {
	grq.ctx.Limit = &limit
	return grq
}

// Offset to start from.
func (grq *GuildRoleQuery) Offset(offset int) *GuildRoleQuery {
	grq.ctx.Offset = &offset
	return grq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (grq *GuildRoleQuery) Unique(unique bool) *GuildRoleQuery {
	grq.ctx.Unique = &unique
	return grq
}

// Order specifies how the records should be ordered.
func (grq *GuildRoleQuery) Order(o ...guildrole.OrderOption) *GuildRoleQuery {
	grq.order = append(grq.order, o...)
	return grq
}

// QueryGuild chains the current query on the "guild" edge.
func (grq *GuildRoleQuery) QueryGuild() *GuildQuery {
	query := (&GuildClient{config: grq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := grq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := grq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(guildrole.Table, guildrole.FieldID, selector),
			sqlgraph.To(guild.Table, guild.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, guildrole.GuildTable, guildrole.GuildColumn),
		)
		fromU = sqlgraph.SetNeighbors(grq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryUserMemberships chains the current query on the "user_memberships" edge.
func (grq *GuildRoleQuery) QueryUserMemberships() *UserMembershipQuery {
	query := (&UserMembershipClient{config: grq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := grq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := grq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(guildrole.Table, guildrole.FieldID, selector),
			sqlgraph.To(usermembership.Table, usermembership.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, guildrole.UserMembershipsTable, guildrole.UserMembershipsPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(grq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryTalent chains the current query on the "talent" edge.
func (grq *GuildRoleQuery) QueryTalent() *YouTubeTalentQuery {
	query := (&YouTubeTalentClient{config: grq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := grq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := grq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(guildrole.Table, guildrole.FieldID, selector),
			sqlgraph.To(youtubetalent.Table, youtubetalent.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, guildrole.TalentTable, guildrole.TalentColumn),
		)
		fromU = sqlgraph.SetNeighbors(grq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first GuildRole entity from the query.
// Returns a *NotFoundError when no GuildRole was found.
func (grq *GuildRoleQuery) First(ctx context.Context) (*GuildRole, error) {
	nodes, err := grq.Limit(1).All(setContextOp(ctx, grq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{guildrole.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (grq *GuildRoleQuery) FirstX(ctx context.Context) *GuildRole {
	node, err := grq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first GuildRole ID from the query.
// Returns a *NotFoundError when no GuildRole ID was found.
func (grq *GuildRoleQuery) FirstID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = grq.Limit(1).IDs(setContextOp(ctx, grq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{guildrole.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (grq *GuildRoleQuery) FirstIDX(ctx context.Context) uint64 {
	id, err := grq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single GuildRole entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one GuildRole entity is found.
// Returns a *NotFoundError when no GuildRole entities are found.
func (grq *GuildRoleQuery) Only(ctx context.Context) (*GuildRole, error) {
	nodes, err := grq.Limit(2).All(setContextOp(ctx, grq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{guildrole.Label}
	default:
		return nil, &NotSingularError{guildrole.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (grq *GuildRoleQuery) OnlyX(ctx context.Context) *GuildRole {
	node, err := grq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only GuildRole ID in the query.
// Returns a *NotSingularError when more than one GuildRole ID is found.
// Returns a *NotFoundError when no entities are found.
func (grq *GuildRoleQuery) OnlyID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = grq.Limit(2).IDs(setContextOp(ctx, grq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{guildrole.Label}
	default:
		err = &NotSingularError{guildrole.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (grq *GuildRoleQuery) OnlyIDX(ctx context.Context) uint64 {
	id, err := grq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of GuildRoles.
func (grq *GuildRoleQuery) All(ctx context.Context) ([]*GuildRole, error) {
	ctx = setContextOp(ctx, grq.ctx, ent.OpQueryAll)
	if err := grq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*GuildRole, *GuildRoleQuery]()
	return withInterceptors[[]*GuildRole](ctx, grq, qr, grq.inters)
}

// AllX is like All, but panics if an error occurs.
func (grq *GuildRoleQuery) AllX(ctx context.Context) []*GuildRole {
	nodes, err := grq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of GuildRole IDs.
func (grq *GuildRoleQuery) IDs(ctx context.Context) (ids []uint64, err error) {
	if grq.ctx.Unique == nil && grq.path != nil {
		grq.Unique(true)
	}
	ctx = setContextOp(ctx, grq.ctx, ent.OpQueryIDs)
	if err = grq.Select(guildrole.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (grq *GuildRoleQuery) IDsX(ctx context.Context) []uint64 {
	ids, err := grq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (grq *GuildRoleQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, grq.ctx, ent.OpQueryCount)
	if err := grq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, grq, querierCount[*GuildRoleQuery](), grq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (grq *GuildRoleQuery) CountX(ctx context.Context) int {
	count, err := grq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (grq *GuildRoleQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, grq.ctx, ent.OpQueryExist)
	switch _, err := grq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (grq *GuildRoleQuery) ExistX(ctx context.Context) bool {
	exist, err := grq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the GuildRoleQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (grq *GuildRoleQuery) Clone() *GuildRoleQuery {
	if grq == nil {
		return nil
	}
	return &GuildRoleQuery{
		config:              grq.config,
		ctx:                 grq.ctx.Clone(),
		order:               append([]guildrole.OrderOption{}, grq.order...),
		inters:              append([]Interceptor{}, grq.inters...),
		predicates:          append([]predicate.GuildRole{}, grq.predicates...),
		withGuild:           grq.withGuild.Clone(),
		withUserMemberships: grq.withUserMemberships.Clone(),
		withTalent:          grq.withTalent.Clone(),
		// clone intermediate query.
		sql:  grq.sql.Clone(),
		path: grq.path,
	}
}

// WithGuild tells the query-builder to eager-load the nodes that are connected to
// the "guild" edge. The optional arguments are used to configure the query builder of the edge.
func (grq *GuildRoleQuery) WithGuild(opts ...func(*GuildQuery)) *GuildRoleQuery {
	query := (&GuildClient{config: grq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	grq.withGuild = query
	return grq
}

// WithUserMemberships tells the query-builder to eager-load the nodes that are connected to
// the "user_memberships" edge. The optional arguments are used to configure the query builder of the edge.
func (grq *GuildRoleQuery) WithUserMemberships(opts ...func(*UserMembershipQuery)) *GuildRoleQuery {
	query := (&UserMembershipClient{config: grq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	grq.withUserMemberships = query
	return grq
}

// WithTalent tells the query-builder to eager-load the nodes that are connected to
// the "talent" edge. The optional arguments are used to configure the query builder of the edge.
func (grq *GuildRoleQuery) WithTalent(opts ...func(*YouTubeTalentQuery)) *GuildRoleQuery {
	query := (&YouTubeTalentClient{config: grq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	grq.withTalent = query
	return grq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.GuildRole.Query().
//		GroupBy(guildrole.FieldName).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (grq *GuildRoleQuery) GroupBy(field string, fields ...string) *GuildRoleGroupBy {
	grq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &GuildRoleGroupBy{build: grq}
	grbuild.flds = &grq.ctx.Fields
	grbuild.label = guildrole.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//	}
//
//	client.GuildRole.Query().
//		Select(guildrole.FieldName).
//		Scan(ctx, &v)
func (grq *GuildRoleQuery) Select(fields ...string) *GuildRoleSelect {
	grq.ctx.Fields = append(grq.ctx.Fields, fields...)
	sbuild := &GuildRoleSelect{GuildRoleQuery: grq}
	sbuild.label = guildrole.Label
	sbuild.flds, sbuild.scan = &grq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a GuildRoleSelect configured with the given aggregations.
func (grq *GuildRoleQuery) Aggregate(fns ...AggregateFunc) *GuildRoleSelect {
	return grq.Select().Aggregate(fns...)
}

func (grq *GuildRoleQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range grq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, grq); err != nil {
				return err
			}
		}
	}
	for _, f := range grq.ctx.Fields {
		if !guildrole.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if grq.path != nil {
		prev, err := grq.path(ctx)
		if err != nil {
			return err
		}
		grq.sql = prev
	}
	return nil
}

func (grq *GuildRoleQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*GuildRole, error) {
	var (
		nodes       = []*GuildRole{}
		withFKs     = grq.withFKs
		_spec       = grq.querySpec()
		loadedTypes = [3]bool{
			grq.withGuild != nil,
			grq.withUserMemberships != nil,
			grq.withTalent != nil,
		}
	)
	if grq.withGuild != nil || grq.withTalent != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, guildrole.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*GuildRole).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &GuildRole{config: grq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(grq.modifiers) > 0 {
		_spec.Modifiers = grq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, grq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := grq.withGuild; query != nil {
		if err := grq.loadGuild(ctx, query, nodes, nil,
			func(n *GuildRole, e *Guild) { n.Edges.Guild = e }); err != nil {
			return nil, err
		}
	}
	if query := grq.withUserMemberships; query != nil {
		if err := grq.loadUserMemberships(ctx, query, nodes,
			func(n *GuildRole) { n.Edges.UserMemberships = []*UserMembership{} },
			func(n *GuildRole, e *UserMembership) { n.Edges.UserMemberships = append(n.Edges.UserMemberships, e) }); err != nil {
			return nil, err
		}
	}
	if query := grq.withTalent; query != nil {
		if err := grq.loadTalent(ctx, query, nodes, nil,
			func(n *GuildRole, e *YouTubeTalent) { n.Edges.Talent = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (grq *GuildRoleQuery) loadGuild(ctx context.Context, query *GuildQuery, nodes []*GuildRole, init func(*GuildRole), assign func(*GuildRole, *Guild)) error {
	ids := make([]uint64, 0, len(nodes))
	nodeids := make(map[uint64][]*GuildRole)
	for i := range nodes {
		if nodes[i].guild_roles == nil {
			continue
		}
		fk := *nodes[i].guild_roles
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(guild.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "guild_roles" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (grq *GuildRoleQuery) loadUserMemberships(ctx context.Context, query *UserMembershipQuery, nodes []*GuildRole, init func(*GuildRole), assign func(*GuildRole, *UserMembership)) error {
	edgeIDs := make([]driver.Value, len(nodes))
	byID := make(map[uint64]*GuildRole)
	nids := make(map[int]map[*GuildRole]struct{})
	for i, node := range nodes {
		edgeIDs[i] = node.ID
		byID[node.ID] = node
		if init != nil {
			init(node)
		}
	}
	query.Where(func(s *sql.Selector) {
		joinT := sql.Table(guildrole.UserMembershipsTable)
		s.Join(joinT).On(s.C(usermembership.FieldID), joinT.C(guildrole.UserMembershipsPrimaryKey[0]))
		s.Where(sql.InValues(joinT.C(guildrole.UserMembershipsPrimaryKey[1]), edgeIDs...))
		columns := s.SelectedColumns()
		s.Select(joinT.C(guildrole.UserMembershipsPrimaryKey[1]))
		s.AppendSelect(columns...)
		s.SetDistinct(false)
	})
	if err := query.prepareQuery(ctx); err != nil {
		return err
	}
	qr := QuerierFunc(func(ctx context.Context, q Query) (Value, error) {
		return query.sqlAll(ctx, func(_ context.Context, spec *sqlgraph.QuerySpec) {
			assign := spec.Assign
			values := spec.ScanValues
			spec.ScanValues = func(columns []string) ([]any, error) {
				values, err := values(columns[1:])
				if err != nil {
					return nil, err
				}
				return append([]any{new(sql.NullInt64)}, values...), nil
			}
			spec.Assign = func(columns []string, values []any) error {
				outValue := uint64(values[0].(*sql.NullInt64).Int64)
				inValue := int(values[1].(*sql.NullInt64).Int64)
				if nids[inValue] == nil {
					nids[inValue] = map[*GuildRole]struct{}{byID[outValue]: {}}
					return assign(columns[1:], values[1:])
				}
				nids[inValue][byID[outValue]] = struct{}{}
				return nil
			}
		})
	})
	neighbors, err := withInterceptors[[]*UserMembership](ctx, query, qr, query.inters)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected "user_memberships" node returned %v`, n.ID)
		}
		for kn := range nodes {
			assign(kn, n)
		}
	}
	return nil
}
func (grq *GuildRoleQuery) loadTalent(ctx context.Context, query *YouTubeTalentQuery, nodes []*GuildRole, init func(*GuildRole), assign func(*GuildRole, *YouTubeTalent)) error {
	ids := make([]string, 0, len(nodes))
	nodeids := make(map[string][]*GuildRole)
	for i := range nodes {
		if nodes[i].you_tube_talent_roles == nil {
			continue
		}
		fk := *nodes[i].you_tube_talent_roles
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(youtubetalent.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "you_tube_talent_roles" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (grq *GuildRoleQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := grq.querySpec()
	if len(grq.modifiers) > 0 {
		_spec.Modifiers = grq.modifiers
	}
	_spec.Node.Columns = grq.ctx.Fields
	if len(grq.ctx.Fields) > 0 {
		_spec.Unique = grq.ctx.Unique != nil && *grq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, grq.driver, _spec)
}

func (grq *GuildRoleQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(guildrole.Table, guildrole.Columns, sqlgraph.NewFieldSpec(guildrole.FieldID, field.TypeUint64))
	_spec.From = grq.sql
	if unique := grq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if grq.path != nil {
		_spec.Unique = true
	}
	if fields := grq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, guildrole.FieldID)
		for i := range fields {
			if fields[i] != guildrole.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := grq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := grq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := grq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := grq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (grq *GuildRoleQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(grq.driver.Dialect())
	t1 := builder.Table(guildrole.Table)
	columns := grq.ctx.Fields
	if len(columns) == 0 {
		columns = guildrole.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if grq.sql != nil {
		selector = grq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if grq.ctx.Unique != nil && *grq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range grq.modifiers {
		m(selector)
	}
	for _, p := range grq.predicates {
		p(selector)
	}
	for _, p := range grq.order {
		p(selector)
	}
	if offset := grq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := grq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ForUpdate locks the selected rows against concurrent updates, and prevent them from being
// updated, deleted or "selected ... for update" by other sessions, until the transaction is
// either committed or rolled-back.
func (grq *GuildRoleQuery) ForUpdate(opts ...sql.LockOption) *GuildRoleQuery {
	if grq.driver.Dialect() == dialect.Postgres {
		grq.Unique(false)
	}
	grq.modifiers = append(grq.modifiers, func(s *sql.Selector) {
		s.ForUpdate(opts...)
	})
	return grq
}

// ForShare behaves similarly to ForUpdate, except that it acquires a shared mode lock
// on any rows that are read. Other sessions can read the rows, but cannot modify them
// until your transaction commits.
func (grq *GuildRoleQuery) ForShare(opts ...sql.LockOption) *GuildRoleQuery {
	if grq.driver.Dialect() == dialect.Postgres {
		grq.Unique(false)
	}
	grq.modifiers = append(grq.modifiers, func(s *sql.Selector) {
		s.ForShare(opts...)
	})
	return grq
}

// GuildRoleGroupBy is the group-by builder for GuildRole entities.
type GuildRoleGroupBy struct {
	selector
	build *GuildRoleQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (grgb *GuildRoleGroupBy) Aggregate(fns ...AggregateFunc) *GuildRoleGroupBy {
	grgb.fns = append(grgb.fns, fns...)
	return grgb
}

// Scan applies the selector query and scans the result into the given value.
func (grgb *GuildRoleGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, grgb.build.ctx, ent.OpQueryGroupBy)
	if err := grgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*GuildRoleQuery, *GuildRoleGroupBy](ctx, grgb.build, grgb, grgb.build.inters, v)
}

func (grgb *GuildRoleGroupBy) sqlScan(ctx context.Context, root *GuildRoleQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(grgb.fns))
	for _, fn := range grgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*grgb.flds)+len(grgb.fns))
		for _, f := range *grgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*grgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := grgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// GuildRoleSelect is the builder for selecting fields of GuildRole entities.
type GuildRoleSelect struct {
	*GuildRoleQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (grs *GuildRoleSelect) Aggregate(fns ...AggregateFunc) *GuildRoleSelect {
	grs.fns = append(grs.fns, fns...)
	return grs
}

// Scan applies the selector query and scans the result into the given value.
func (grs *GuildRoleSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, grs.ctx, ent.OpQuerySelect)
	if err := grs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*GuildRoleQuery, *GuildRoleSelect](ctx, grs.GuildRoleQuery, grs, grs.inters, v)
}

func (grs *GuildRoleSelect) sqlScan(ctx context.Context, root *GuildRoleQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(grs.fns))
	for _, fn := range grs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*grs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := grs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
