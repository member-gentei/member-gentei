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
	"github.com/member-gentei/member-gentei/gentei/ent/guildrole"
	"github.com/member-gentei/member-gentei/gentei/ent/predicate"
	"github.com/member-gentei/member-gentei/gentei/ent/user"
	"github.com/member-gentei/member-gentei/gentei/ent/usermembership"
	"github.com/member-gentei/member-gentei/gentei/ent/youtubetalent"
)

// UserMembershipQuery is the builder for querying UserMembership entities.
type UserMembershipQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.UserMembership
	// eager-loading edges.
	withUser          *UserQuery
	withYoutubeTalent *YouTubeTalentQuery
	withRoles         *GuildRoleQuery
	withFKs           bool
	modifiers         []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the UserMembershipQuery builder.
func (umq *UserMembershipQuery) Where(ps ...predicate.UserMembership) *UserMembershipQuery {
	umq.predicates = append(umq.predicates, ps...)
	return umq
}

// Limit adds a limit step to the query.
func (umq *UserMembershipQuery) Limit(limit int) *UserMembershipQuery {
	umq.limit = &limit
	return umq
}

// Offset adds an offset step to the query.
func (umq *UserMembershipQuery) Offset(offset int) *UserMembershipQuery {
	umq.offset = &offset
	return umq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (umq *UserMembershipQuery) Unique(unique bool) *UserMembershipQuery {
	umq.unique = &unique
	return umq
}

// Order adds an order step to the query.
func (umq *UserMembershipQuery) Order(o ...OrderFunc) *UserMembershipQuery {
	umq.order = append(umq.order, o...)
	return umq
}

// QueryUser chains the current query on the "user" edge.
func (umq *UserMembershipQuery) QueryUser() *UserQuery {
	query := &UserQuery{config: umq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := umq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := umq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(usermembership.Table, usermembership.FieldID, selector),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, usermembership.UserTable, usermembership.UserColumn),
		)
		fromU = sqlgraph.SetNeighbors(umq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryYoutubeTalent chains the current query on the "youtube_talent" edge.
func (umq *UserMembershipQuery) QueryYoutubeTalent() *YouTubeTalentQuery {
	query := &YouTubeTalentQuery{config: umq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := umq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := umq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(usermembership.Table, usermembership.FieldID, selector),
			sqlgraph.To(youtubetalent.Table, youtubetalent.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, usermembership.YoutubeTalentTable, usermembership.YoutubeTalentColumn),
		)
		fromU = sqlgraph.SetNeighbors(umq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryRoles chains the current query on the "roles" edge.
func (umq *UserMembershipQuery) QueryRoles() *GuildRoleQuery {
	query := &GuildRoleQuery{config: umq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := umq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := umq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(usermembership.Table, usermembership.FieldID, selector),
			sqlgraph.To(guildrole.Table, guildrole.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, usermembership.RolesTable, usermembership.RolesPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(umq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first UserMembership entity from the query.
// Returns a *NotFoundError when no UserMembership was found.
func (umq *UserMembershipQuery) First(ctx context.Context) (*UserMembership, error) {
	nodes, err := umq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{usermembership.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (umq *UserMembershipQuery) FirstX(ctx context.Context) *UserMembership {
	node, err := umq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first UserMembership ID from the query.
// Returns a *NotFoundError when no UserMembership ID was found.
func (umq *UserMembershipQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = umq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{usermembership.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (umq *UserMembershipQuery) FirstIDX(ctx context.Context) int {
	id, err := umq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single UserMembership entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one UserMembership entity is found.
// Returns a *NotFoundError when no UserMembership entities are found.
func (umq *UserMembershipQuery) Only(ctx context.Context) (*UserMembership, error) {
	nodes, err := umq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{usermembership.Label}
	default:
		return nil, &NotSingularError{usermembership.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (umq *UserMembershipQuery) OnlyX(ctx context.Context) *UserMembership {
	node, err := umq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only UserMembership ID in the query.
// Returns a *NotSingularError when more than one UserMembership ID is found.
// Returns a *NotFoundError when no entities are found.
func (umq *UserMembershipQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = umq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{usermembership.Label}
	default:
		err = &NotSingularError{usermembership.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (umq *UserMembershipQuery) OnlyIDX(ctx context.Context) int {
	id, err := umq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of UserMemberships.
func (umq *UserMembershipQuery) All(ctx context.Context) ([]*UserMembership, error) {
	if err := umq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return umq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (umq *UserMembershipQuery) AllX(ctx context.Context) []*UserMembership {
	nodes, err := umq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of UserMembership IDs.
func (umq *UserMembershipQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := umq.Select(usermembership.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (umq *UserMembershipQuery) IDsX(ctx context.Context) []int {
	ids, err := umq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (umq *UserMembershipQuery) Count(ctx context.Context) (int, error) {
	if err := umq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return umq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (umq *UserMembershipQuery) CountX(ctx context.Context) int {
	count, err := umq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (umq *UserMembershipQuery) Exist(ctx context.Context) (bool, error) {
	if err := umq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return umq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (umq *UserMembershipQuery) ExistX(ctx context.Context) bool {
	exist, err := umq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the UserMembershipQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (umq *UserMembershipQuery) Clone() *UserMembershipQuery {
	if umq == nil {
		return nil
	}
	return &UserMembershipQuery{
		config:            umq.config,
		limit:             umq.limit,
		offset:            umq.offset,
		order:             append([]OrderFunc{}, umq.order...),
		predicates:        append([]predicate.UserMembership{}, umq.predicates...),
		withUser:          umq.withUser.Clone(),
		withYoutubeTalent: umq.withYoutubeTalent.Clone(),
		withRoles:         umq.withRoles.Clone(),
		// clone intermediate query.
		sql:    umq.sql.Clone(),
		path:   umq.path,
		unique: umq.unique,
	}
}

// WithUser tells the query-builder to eager-load the nodes that are connected to
// the "user" edge. The optional arguments are used to configure the query builder of the edge.
func (umq *UserMembershipQuery) WithUser(opts ...func(*UserQuery)) *UserMembershipQuery {
	query := &UserQuery{config: umq.config}
	for _, opt := range opts {
		opt(query)
	}
	umq.withUser = query
	return umq
}

// WithYoutubeTalent tells the query-builder to eager-load the nodes that are connected to
// the "youtube_talent" edge. The optional arguments are used to configure the query builder of the edge.
func (umq *UserMembershipQuery) WithYoutubeTalent(opts ...func(*YouTubeTalentQuery)) *UserMembershipQuery {
	query := &YouTubeTalentQuery{config: umq.config}
	for _, opt := range opts {
		opt(query)
	}
	umq.withYoutubeTalent = query
	return umq
}

// WithRoles tells the query-builder to eager-load the nodes that are connected to
// the "roles" edge. The optional arguments are used to configure the query builder of the edge.
func (umq *UserMembershipQuery) WithRoles(opts ...func(*GuildRoleQuery)) *UserMembershipQuery {
	query := &GuildRoleQuery{config: umq.config}
	for _, opt := range opts {
		opt(query)
	}
	umq.withRoles = query
	return umq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		FirstFailed time.Time `json:"first_failed,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.UserMembership.Query().
//		GroupBy(usermembership.FieldFirstFailed).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (umq *UserMembershipQuery) GroupBy(field string, fields ...string) *UserMembershipGroupBy {
	grbuild := &UserMembershipGroupBy{config: umq.config}
	grbuild.fields = append([]string{field}, fields...)
	grbuild.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := umq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return umq.sqlQuery(ctx), nil
	}
	grbuild.label = usermembership.Label
	grbuild.flds, grbuild.scan = &grbuild.fields, grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		FirstFailed time.Time `json:"first_failed,omitempty"`
//	}
//
//	client.UserMembership.Query().
//		Select(usermembership.FieldFirstFailed).
//		Scan(ctx, &v)
//
func (umq *UserMembershipQuery) Select(fields ...string) *UserMembershipSelect {
	umq.fields = append(umq.fields, fields...)
	selbuild := &UserMembershipSelect{UserMembershipQuery: umq}
	selbuild.label = usermembership.Label
	selbuild.flds, selbuild.scan = &umq.fields, selbuild.Scan
	return selbuild
}

func (umq *UserMembershipQuery) prepareQuery(ctx context.Context) error {
	for _, f := range umq.fields {
		if !usermembership.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if umq.path != nil {
		prev, err := umq.path(ctx)
		if err != nil {
			return err
		}
		umq.sql = prev
	}
	return nil
}

func (umq *UserMembershipQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*UserMembership, error) {
	var (
		nodes       = []*UserMembership{}
		withFKs     = umq.withFKs
		_spec       = umq.querySpec()
		loadedTypes = [3]bool{
			umq.withUser != nil,
			umq.withYoutubeTalent != nil,
			umq.withRoles != nil,
		}
	)
	if umq.withUser != nil || umq.withYoutubeTalent != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, usermembership.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]interface{}, error) {
		return (*UserMembership).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []interface{}) error {
		node := &UserMembership{config: umq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(umq.modifiers) > 0 {
		_spec.Modifiers = umq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, umq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := umq.withUser; query != nil {
		ids := make([]uint64, 0, len(nodes))
		nodeids := make(map[uint64][]*UserMembership)
		for i := range nodes {
			if nodes[i].user_memberships == nil {
				continue
			}
			fk := *nodes[i].user_memberships
			if _, ok := nodeids[fk]; !ok {
				ids = append(ids, fk)
			}
			nodeids[fk] = append(nodeids[fk], nodes[i])
		}
		query.Where(user.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "user_memberships" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.User = n
			}
		}
	}

	if query := umq.withYoutubeTalent; query != nil {
		ids := make([]string, 0, len(nodes))
		nodeids := make(map[string][]*UserMembership)
		for i := range nodes {
			if nodes[i].user_membership_youtube_talent == nil {
				continue
			}
			fk := *nodes[i].user_membership_youtube_talent
			if _, ok := nodeids[fk]; !ok {
				ids = append(ids, fk)
			}
			nodeids[fk] = append(nodeids[fk], nodes[i])
		}
		query.Where(youtubetalent.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "user_membership_youtube_talent" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.YoutubeTalent = n
			}
		}
	}

	if query := umq.withRoles; query != nil {
		edgeids := make([]driver.Value, len(nodes))
		byid := make(map[int]*UserMembership)
		nids := make(map[uint64]map[*UserMembership]struct{})
		for i, node := range nodes {
			edgeids[i] = node.ID
			byid[node.ID] = node
			node.Edges.Roles = []*GuildRole{}
		}
		query.Where(func(s *sql.Selector) {
			joinT := sql.Table(usermembership.RolesTable)
			s.Join(joinT).On(s.C(guildrole.FieldID), joinT.C(usermembership.RolesPrimaryKey[1]))
			s.Where(sql.InValues(joinT.C(usermembership.RolesPrimaryKey[0]), edgeids...))
			columns := s.SelectedColumns()
			s.Select(joinT.C(usermembership.RolesPrimaryKey[0]))
			s.AppendSelect(columns...)
			s.SetDistinct(false)
		})
		neighbors, err := query.sqlAll(ctx, func(_ context.Context, spec *sqlgraph.QuerySpec) {
			assign := spec.Assign
			values := spec.ScanValues
			spec.ScanValues = func(columns []string) ([]interface{}, error) {
				values, err := values(columns[1:])
				if err != nil {
					return nil, err
				}
				return append([]interface{}{new(sql.NullInt64)}, values...), nil
			}
			spec.Assign = func(columns []string, values []interface{}) error {
				outValue := int(values[0].(*sql.NullInt64).Int64)
				inValue := uint64(values[1].(*sql.NullInt64).Int64)
				if nids[inValue] == nil {
					nids[inValue] = map[*UserMembership]struct{}{byid[outValue]: struct{}{}}
					return assign(columns[1:], values[1:])
				}
				nids[inValue][byid[outValue]] = struct{}{}
				return nil
			}
		})
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected "roles" node returned %v`, n.ID)
			}
			for kn := range nodes {
				kn.Edges.Roles = append(kn.Edges.Roles, n)
			}
		}
	}

	return nodes, nil
}

func (umq *UserMembershipQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := umq.querySpec()
	if len(umq.modifiers) > 0 {
		_spec.Modifiers = umq.modifiers
	}
	_spec.Node.Columns = umq.fields
	if len(umq.fields) > 0 {
		_spec.Unique = umq.unique != nil && *umq.unique
	}
	return sqlgraph.CountNodes(ctx, umq.driver, _spec)
}

func (umq *UserMembershipQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := umq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (umq *UserMembershipQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   usermembership.Table,
			Columns: usermembership.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: usermembership.FieldID,
			},
		},
		From:   umq.sql,
		Unique: true,
	}
	if unique := umq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := umq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, usermembership.FieldID)
		for i := range fields {
			if fields[i] != usermembership.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := umq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := umq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := umq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := umq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (umq *UserMembershipQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(umq.driver.Dialect())
	t1 := builder.Table(usermembership.Table)
	columns := umq.fields
	if len(columns) == 0 {
		columns = usermembership.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if umq.sql != nil {
		selector = umq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if umq.unique != nil && *umq.unique {
		selector.Distinct()
	}
	for _, m := range umq.modifiers {
		m(selector)
	}
	for _, p := range umq.predicates {
		p(selector)
	}
	for _, p := range umq.order {
		p(selector)
	}
	if offset := umq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := umq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ForUpdate locks the selected rows against concurrent updates, and prevent them from being
// updated, deleted or "selected ... for update" by other sessions, until the transaction is
// either committed or rolled-back.
func (umq *UserMembershipQuery) ForUpdate(opts ...sql.LockOption) *UserMembershipQuery {
	if umq.driver.Dialect() == dialect.Postgres {
		umq.Unique(false)
	}
	umq.modifiers = append(umq.modifiers, func(s *sql.Selector) {
		s.ForUpdate(opts...)
	})
	return umq
}

// ForShare behaves similarly to ForUpdate, except that it acquires a shared mode lock
// on any rows that are read. Other sessions can read the rows, but cannot modify them
// until your transaction commits.
func (umq *UserMembershipQuery) ForShare(opts ...sql.LockOption) *UserMembershipQuery {
	if umq.driver.Dialect() == dialect.Postgres {
		umq.Unique(false)
	}
	umq.modifiers = append(umq.modifiers, func(s *sql.Selector) {
		s.ForShare(opts...)
	})
	return umq
}

// UserMembershipGroupBy is the group-by builder for UserMembership entities.
type UserMembershipGroupBy struct {
	config
	selector
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (umgb *UserMembershipGroupBy) Aggregate(fns ...AggregateFunc) *UserMembershipGroupBy {
	umgb.fns = append(umgb.fns, fns...)
	return umgb
}

// Scan applies the group-by query and scans the result into the given value.
func (umgb *UserMembershipGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := umgb.path(ctx)
	if err != nil {
		return err
	}
	umgb.sql = query
	return umgb.sqlScan(ctx, v)
}

func (umgb *UserMembershipGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	for _, f := range umgb.fields {
		if !usermembership.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := umgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := umgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (umgb *UserMembershipGroupBy) sqlQuery() *sql.Selector {
	selector := umgb.sql.Select()
	aggregation := make([]string, 0, len(umgb.fns))
	for _, fn := range umgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(umgb.fields)+len(umgb.fns))
		for _, f := range umgb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(umgb.fields...)...)
}

// UserMembershipSelect is the builder for selecting fields of UserMembership entities.
type UserMembershipSelect struct {
	*UserMembershipQuery
	selector
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (ums *UserMembershipSelect) Scan(ctx context.Context, v interface{}) error {
	if err := ums.prepareQuery(ctx); err != nil {
		return err
	}
	ums.sql = ums.UserMembershipQuery.sqlQuery(ctx)
	return ums.sqlScan(ctx, v)
}

func (ums *UserMembershipSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := ums.sql.Query()
	if err := ums.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
