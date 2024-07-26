// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

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

// YouTubeTalentQuery is the builder for querying YouTubeTalent entities.
type YouTubeTalentQuery struct {
	config
	ctx             *QueryContext
	order           []youtubetalent.OrderOption
	inters          []Interceptor
	predicates      []predicate.YouTubeTalent
	withGuilds      *GuildQuery
	withRoles       *GuildRoleQuery
	withMemberships *UserMembershipQuery
	modifiers       []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the YouTubeTalentQuery builder.
func (yttq *YouTubeTalentQuery) Where(ps ...predicate.YouTubeTalent) *YouTubeTalentQuery {
	yttq.predicates = append(yttq.predicates, ps...)
	return yttq
}

// Limit the number of records to be returned by this query.
func (yttq *YouTubeTalentQuery) Limit(limit int) *YouTubeTalentQuery {
	yttq.ctx.Limit = &limit
	return yttq
}

// Offset to start from.
func (yttq *YouTubeTalentQuery) Offset(offset int) *YouTubeTalentQuery {
	yttq.ctx.Offset = &offset
	return yttq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (yttq *YouTubeTalentQuery) Unique(unique bool) *YouTubeTalentQuery {
	yttq.ctx.Unique = &unique
	return yttq
}

// Order specifies how the records should be ordered.
func (yttq *YouTubeTalentQuery) Order(o ...youtubetalent.OrderOption) *YouTubeTalentQuery {
	yttq.order = append(yttq.order, o...)
	return yttq
}

// QueryGuilds chains the current query on the "guilds" edge.
func (yttq *YouTubeTalentQuery) QueryGuilds() *GuildQuery {
	query := (&GuildClient{config: yttq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := yttq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := yttq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(youtubetalent.Table, youtubetalent.FieldID, selector),
			sqlgraph.To(guild.Table, guild.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, youtubetalent.GuildsTable, youtubetalent.GuildsPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(yttq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryRoles chains the current query on the "roles" edge.
func (yttq *YouTubeTalentQuery) QueryRoles() *GuildRoleQuery {
	query := (&GuildRoleClient{config: yttq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := yttq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := yttq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(youtubetalent.Table, youtubetalent.FieldID, selector),
			sqlgraph.To(guildrole.Table, guildrole.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, youtubetalent.RolesTable, youtubetalent.RolesColumn),
		)
		fromU = sqlgraph.SetNeighbors(yttq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryMemberships chains the current query on the "memberships" edge.
func (yttq *YouTubeTalentQuery) QueryMemberships() *UserMembershipQuery {
	query := (&UserMembershipClient{config: yttq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := yttq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := yttq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(youtubetalent.Table, youtubetalent.FieldID, selector),
			sqlgraph.To(usermembership.Table, usermembership.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, youtubetalent.MembershipsTable, youtubetalent.MembershipsColumn),
		)
		fromU = sqlgraph.SetNeighbors(yttq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first YouTubeTalent entity from the query.
// Returns a *NotFoundError when no YouTubeTalent was found.
func (yttq *YouTubeTalentQuery) First(ctx context.Context) (*YouTubeTalent, error) {
	nodes, err := yttq.Limit(1).All(setContextOp(ctx, yttq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{youtubetalent.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (yttq *YouTubeTalentQuery) FirstX(ctx context.Context) *YouTubeTalent {
	node, err := yttq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first YouTubeTalent ID from the query.
// Returns a *NotFoundError when no YouTubeTalent ID was found.
func (yttq *YouTubeTalentQuery) FirstID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = yttq.Limit(1).IDs(setContextOp(ctx, yttq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{youtubetalent.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (yttq *YouTubeTalentQuery) FirstIDX(ctx context.Context) string {
	id, err := yttq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single YouTubeTalent entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one YouTubeTalent entity is found.
// Returns a *NotFoundError when no YouTubeTalent entities are found.
func (yttq *YouTubeTalentQuery) Only(ctx context.Context) (*YouTubeTalent, error) {
	nodes, err := yttq.Limit(2).All(setContextOp(ctx, yttq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{youtubetalent.Label}
	default:
		return nil, &NotSingularError{youtubetalent.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (yttq *YouTubeTalentQuery) OnlyX(ctx context.Context) *YouTubeTalent {
	node, err := yttq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only YouTubeTalent ID in the query.
// Returns a *NotSingularError when more than one YouTubeTalent ID is found.
// Returns a *NotFoundError when no entities are found.
func (yttq *YouTubeTalentQuery) OnlyID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = yttq.Limit(2).IDs(setContextOp(ctx, yttq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{youtubetalent.Label}
	default:
		err = &NotSingularError{youtubetalent.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (yttq *YouTubeTalentQuery) OnlyIDX(ctx context.Context) string {
	id, err := yttq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of YouTubeTalents.
func (yttq *YouTubeTalentQuery) All(ctx context.Context) ([]*YouTubeTalent, error) {
	ctx = setContextOp(ctx, yttq.ctx, "All")
	if err := yttq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*YouTubeTalent, *YouTubeTalentQuery]()
	return withInterceptors[[]*YouTubeTalent](ctx, yttq, qr, yttq.inters)
}

// AllX is like All, but panics if an error occurs.
func (yttq *YouTubeTalentQuery) AllX(ctx context.Context) []*YouTubeTalent {
	nodes, err := yttq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of YouTubeTalent IDs.
func (yttq *YouTubeTalentQuery) IDs(ctx context.Context) (ids []string, err error) {
	if yttq.ctx.Unique == nil && yttq.path != nil {
		yttq.Unique(true)
	}
	ctx = setContextOp(ctx, yttq.ctx, "IDs")
	if err = yttq.Select(youtubetalent.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (yttq *YouTubeTalentQuery) IDsX(ctx context.Context) []string {
	ids, err := yttq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (yttq *YouTubeTalentQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, yttq.ctx, "Count")
	if err := yttq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, yttq, querierCount[*YouTubeTalentQuery](), yttq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (yttq *YouTubeTalentQuery) CountX(ctx context.Context) int {
	count, err := yttq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (yttq *YouTubeTalentQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, yttq.ctx, "Exist")
	switch _, err := yttq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (yttq *YouTubeTalentQuery) ExistX(ctx context.Context) bool {
	exist, err := yttq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the YouTubeTalentQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (yttq *YouTubeTalentQuery) Clone() *YouTubeTalentQuery {
	if yttq == nil {
		return nil
	}
	return &YouTubeTalentQuery{
		config:          yttq.config,
		ctx:             yttq.ctx.Clone(),
		order:           append([]youtubetalent.OrderOption{}, yttq.order...),
		inters:          append([]Interceptor{}, yttq.inters...),
		predicates:      append([]predicate.YouTubeTalent{}, yttq.predicates...),
		withGuilds:      yttq.withGuilds.Clone(),
		withRoles:       yttq.withRoles.Clone(),
		withMemberships: yttq.withMemberships.Clone(),
		// clone intermediate query.
		sql:  yttq.sql.Clone(),
		path: yttq.path,
	}
}

// WithGuilds tells the query-builder to eager-load the nodes that are connected to
// the "guilds" edge. The optional arguments are used to configure the query builder of the edge.
func (yttq *YouTubeTalentQuery) WithGuilds(opts ...func(*GuildQuery)) *YouTubeTalentQuery {
	query := (&GuildClient{config: yttq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	yttq.withGuilds = query
	return yttq
}

// WithRoles tells the query-builder to eager-load the nodes that are connected to
// the "roles" edge. The optional arguments are used to configure the query builder of the edge.
func (yttq *YouTubeTalentQuery) WithRoles(opts ...func(*GuildRoleQuery)) *YouTubeTalentQuery {
	query := (&GuildRoleClient{config: yttq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	yttq.withRoles = query
	return yttq
}

// WithMemberships tells the query-builder to eager-load the nodes that are connected to
// the "memberships" edge. The optional arguments are used to configure the query builder of the edge.
func (yttq *YouTubeTalentQuery) WithMemberships(opts ...func(*UserMembershipQuery)) *YouTubeTalentQuery {
	query := (&UserMembershipClient{config: yttq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	yttq.withMemberships = query
	return yttq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		ChannelName string `json:"channel_name,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.YouTubeTalent.Query().
//		GroupBy(youtubetalent.FieldChannelName).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (yttq *YouTubeTalentQuery) GroupBy(field string, fields ...string) *YouTubeTalentGroupBy {
	yttq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &YouTubeTalentGroupBy{build: yttq}
	grbuild.flds = &yttq.ctx.Fields
	grbuild.label = youtubetalent.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		ChannelName string `json:"channel_name,omitempty"`
//	}
//
//	client.YouTubeTalent.Query().
//		Select(youtubetalent.FieldChannelName).
//		Scan(ctx, &v)
func (yttq *YouTubeTalentQuery) Select(fields ...string) *YouTubeTalentSelect {
	yttq.ctx.Fields = append(yttq.ctx.Fields, fields...)
	sbuild := &YouTubeTalentSelect{YouTubeTalentQuery: yttq}
	sbuild.label = youtubetalent.Label
	sbuild.flds, sbuild.scan = &yttq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a YouTubeTalentSelect configured with the given aggregations.
func (yttq *YouTubeTalentQuery) Aggregate(fns ...AggregateFunc) *YouTubeTalentSelect {
	return yttq.Select().Aggregate(fns...)
}

func (yttq *YouTubeTalentQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range yttq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, yttq); err != nil {
				return err
			}
		}
	}
	for _, f := range yttq.ctx.Fields {
		if !youtubetalent.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if yttq.path != nil {
		prev, err := yttq.path(ctx)
		if err != nil {
			return err
		}
		yttq.sql = prev
	}
	return nil
}

func (yttq *YouTubeTalentQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*YouTubeTalent, error) {
	var (
		nodes       = []*YouTubeTalent{}
		_spec       = yttq.querySpec()
		loadedTypes = [3]bool{
			yttq.withGuilds != nil,
			yttq.withRoles != nil,
			yttq.withMemberships != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*YouTubeTalent).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &YouTubeTalent{config: yttq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(yttq.modifiers) > 0 {
		_spec.Modifiers = yttq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, yttq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := yttq.withGuilds; query != nil {
		if err := yttq.loadGuilds(ctx, query, nodes,
			func(n *YouTubeTalent) { n.Edges.Guilds = []*Guild{} },
			func(n *YouTubeTalent, e *Guild) { n.Edges.Guilds = append(n.Edges.Guilds, e) }); err != nil {
			return nil, err
		}
	}
	if query := yttq.withRoles; query != nil {
		if err := yttq.loadRoles(ctx, query, nodes,
			func(n *YouTubeTalent) { n.Edges.Roles = []*GuildRole{} },
			func(n *YouTubeTalent, e *GuildRole) { n.Edges.Roles = append(n.Edges.Roles, e) }); err != nil {
			return nil, err
		}
	}
	if query := yttq.withMemberships; query != nil {
		if err := yttq.loadMemberships(ctx, query, nodes,
			func(n *YouTubeTalent) { n.Edges.Memberships = []*UserMembership{} },
			func(n *YouTubeTalent, e *UserMembership) { n.Edges.Memberships = append(n.Edges.Memberships, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (yttq *YouTubeTalentQuery) loadGuilds(ctx context.Context, query *GuildQuery, nodes []*YouTubeTalent, init func(*YouTubeTalent), assign func(*YouTubeTalent, *Guild)) error {
	edgeIDs := make([]driver.Value, len(nodes))
	byID := make(map[string]*YouTubeTalent)
	nids := make(map[uint64]map[*YouTubeTalent]struct{})
	for i, node := range nodes {
		edgeIDs[i] = node.ID
		byID[node.ID] = node
		if init != nil {
			init(node)
		}
	}
	query.Where(func(s *sql.Selector) {
		joinT := sql.Table(youtubetalent.GuildsTable)
		s.Join(joinT).On(s.C(guild.FieldID), joinT.C(youtubetalent.GuildsPrimaryKey[1]))
		s.Where(sql.InValues(joinT.C(youtubetalent.GuildsPrimaryKey[0]), edgeIDs...))
		columns := s.SelectedColumns()
		s.Select(joinT.C(youtubetalent.GuildsPrimaryKey[0]))
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
				return append([]any{new(sql.NullString)}, values...), nil
			}
			spec.Assign = func(columns []string, values []any) error {
				outValue := values[0].(*sql.NullString).String
				inValue := uint64(values[1].(*sql.NullInt64).Int64)
				if nids[inValue] == nil {
					nids[inValue] = map[*YouTubeTalent]struct{}{byID[outValue]: {}}
					return assign(columns[1:], values[1:])
				}
				nids[inValue][byID[outValue]] = struct{}{}
				return nil
			}
		})
	})
	neighbors, err := withInterceptors[[]*Guild](ctx, query, qr, query.inters)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected "guilds" node returned %v`, n.ID)
		}
		for kn := range nodes {
			assign(kn, n)
		}
	}
	return nil
}
func (yttq *YouTubeTalentQuery) loadRoles(ctx context.Context, query *GuildRoleQuery, nodes []*YouTubeTalent, init func(*YouTubeTalent), assign func(*YouTubeTalent, *GuildRole)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[string]*YouTubeTalent)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.GuildRole(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(youtubetalent.RolesColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.you_tube_talent_roles
		if fk == nil {
			return fmt.Errorf(`foreign-key "you_tube_talent_roles" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "you_tube_talent_roles" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (yttq *YouTubeTalentQuery) loadMemberships(ctx context.Context, query *UserMembershipQuery, nodes []*YouTubeTalent, init func(*YouTubeTalent), assign func(*YouTubeTalent, *UserMembership)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[string]*YouTubeTalent)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.UserMembership(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(youtubetalent.MembershipsColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.user_membership_youtube_talent
		if fk == nil {
			return fmt.Errorf(`foreign-key "user_membership_youtube_talent" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "user_membership_youtube_talent" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (yttq *YouTubeTalentQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := yttq.querySpec()
	if len(yttq.modifiers) > 0 {
		_spec.Modifiers = yttq.modifiers
	}
	_spec.Node.Columns = yttq.ctx.Fields
	if len(yttq.ctx.Fields) > 0 {
		_spec.Unique = yttq.ctx.Unique != nil && *yttq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, yttq.driver, _spec)
}

func (yttq *YouTubeTalentQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(youtubetalent.Table, youtubetalent.Columns, sqlgraph.NewFieldSpec(youtubetalent.FieldID, field.TypeString))
	_spec.From = yttq.sql
	if unique := yttq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if yttq.path != nil {
		_spec.Unique = true
	}
	if fields := yttq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, youtubetalent.FieldID)
		for i := range fields {
			if fields[i] != youtubetalent.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := yttq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := yttq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := yttq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := yttq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (yttq *YouTubeTalentQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(yttq.driver.Dialect())
	t1 := builder.Table(youtubetalent.Table)
	columns := yttq.ctx.Fields
	if len(columns) == 0 {
		columns = youtubetalent.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if yttq.sql != nil {
		selector = yttq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if yttq.ctx.Unique != nil && *yttq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range yttq.modifiers {
		m(selector)
	}
	for _, p := range yttq.predicates {
		p(selector)
	}
	for _, p := range yttq.order {
		p(selector)
	}
	if offset := yttq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := yttq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ForUpdate locks the selected rows against concurrent updates, and prevent them from being
// updated, deleted or "selected ... for update" by other sessions, until the transaction is
// either committed or rolled-back.
func (yttq *YouTubeTalentQuery) ForUpdate(opts ...sql.LockOption) *YouTubeTalentQuery {
	if yttq.driver.Dialect() == dialect.Postgres {
		yttq.Unique(false)
	}
	yttq.modifiers = append(yttq.modifiers, func(s *sql.Selector) {
		s.ForUpdate(opts...)
	})
	return yttq
}

// ForShare behaves similarly to ForUpdate, except that it acquires a shared mode lock
// on any rows that are read. Other sessions can read the rows, but cannot modify them
// until your transaction commits.
func (yttq *YouTubeTalentQuery) ForShare(opts ...sql.LockOption) *YouTubeTalentQuery {
	if yttq.driver.Dialect() == dialect.Postgres {
		yttq.Unique(false)
	}
	yttq.modifiers = append(yttq.modifiers, func(s *sql.Selector) {
		s.ForShare(opts...)
	})
	return yttq
}

// YouTubeTalentGroupBy is the group-by builder for YouTubeTalent entities.
type YouTubeTalentGroupBy struct {
	selector
	build *YouTubeTalentQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (yttgb *YouTubeTalentGroupBy) Aggregate(fns ...AggregateFunc) *YouTubeTalentGroupBy {
	yttgb.fns = append(yttgb.fns, fns...)
	return yttgb
}

// Scan applies the selector query and scans the result into the given value.
func (yttgb *YouTubeTalentGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, yttgb.build.ctx, "GroupBy")
	if err := yttgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*YouTubeTalentQuery, *YouTubeTalentGroupBy](ctx, yttgb.build, yttgb, yttgb.build.inters, v)
}

func (yttgb *YouTubeTalentGroupBy) sqlScan(ctx context.Context, root *YouTubeTalentQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(yttgb.fns))
	for _, fn := range yttgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*yttgb.flds)+len(yttgb.fns))
		for _, f := range *yttgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*yttgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := yttgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// YouTubeTalentSelect is the builder for selecting fields of YouTubeTalent entities.
type YouTubeTalentSelect struct {
	*YouTubeTalentQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ytts *YouTubeTalentSelect) Aggregate(fns ...AggregateFunc) *YouTubeTalentSelect {
	ytts.fns = append(ytts.fns, fns...)
	return ytts
}

// Scan applies the selector query and scans the result into the given value.
func (ytts *YouTubeTalentSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ytts.ctx, "Select")
	if err := ytts.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*YouTubeTalentQuery, *YouTubeTalentSelect](ctx, ytts.YouTubeTalentQuery, ytts, ytts.inters, v)
}

func (ytts *YouTubeTalentSelect) sqlScan(ctx context.Context, root *YouTubeTalentQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ytts.fns))
	for _, fn := range ytts.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ytts.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ytts.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
