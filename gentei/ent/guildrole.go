// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/member-gentei/member-gentei/gentei/ent/guild"
	"github.com/member-gentei/member-gentei/gentei/ent/guildrole"
	"github.com/member-gentei/member-gentei/gentei/ent/youtubetalent"
)

// GuildRole is the model entity for the GuildRole schema.
type GuildRole struct {
	config `json:"-"`
	// ID of the ent.
	// Discord snowflake for this role
	ID uint64 `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	// Human name for this role
	Name string `json:"name,omitempty"`
	// LastUpdated holds the value of the "last_updated" field.
	// When the name was last synchronized for this role
	LastUpdated time.Time `json:"last_updated,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the GuildRoleQuery when eager-loading is set.
	Edges                 GuildRoleEdges `json:"edges"`
	guild_roles           *uint64
	you_tube_talent_roles *string
}

// GuildRoleEdges holds the relations/edges for other nodes in the graph.
type GuildRoleEdges struct {
	// Guild holds the value of the guild edge.
	Guild *Guild `json:"guild,omitempty"`
	// UserMemberships holds the value of the user_memberships edge.
	UserMemberships []*UserMembership `json:"user_memberships,omitempty"`
	// Talent holds the value of the talent edge.
	Talent *YouTubeTalent `json:"talent,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// GuildOrErr returns the Guild value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e GuildRoleEdges) GuildOrErr() (*Guild, error) {
	if e.loadedTypes[0] {
		if e.Guild == nil {
			// The edge guild was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: guild.Label}
		}
		return e.Guild, nil
	}
	return nil, &NotLoadedError{edge: "guild"}
}

// UserMembershipsOrErr returns the UserMemberships value or an error if the edge
// was not loaded in eager-loading.
func (e GuildRoleEdges) UserMembershipsOrErr() ([]*UserMembership, error) {
	if e.loadedTypes[1] {
		return e.UserMemberships, nil
	}
	return nil, &NotLoadedError{edge: "user_memberships"}
}

// TalentOrErr returns the Talent value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e GuildRoleEdges) TalentOrErr() (*YouTubeTalent, error) {
	if e.loadedTypes[2] {
		if e.Talent == nil {
			// The edge talent was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: youtubetalent.Label}
		}
		return e.Talent, nil
	}
	return nil, &NotLoadedError{edge: "talent"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*GuildRole) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case guildrole.FieldID:
			values[i] = new(sql.NullInt64)
		case guildrole.FieldName:
			values[i] = new(sql.NullString)
		case guildrole.FieldLastUpdated:
			values[i] = new(sql.NullTime)
		case guildrole.ForeignKeys[0]: // guild_roles
			values[i] = new(sql.NullInt64)
		case guildrole.ForeignKeys[1]: // you_tube_talent_roles
			values[i] = new(sql.NullString)
		default:
			return nil, fmt.Errorf("unexpected column %q for type GuildRole", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the GuildRole fields.
func (gr *GuildRole) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case guildrole.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			gr.ID = uint64(value.Int64)
		case guildrole.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				gr.Name = value.String
			}
		case guildrole.FieldLastUpdated:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field last_updated", values[i])
			} else if value.Valid {
				gr.LastUpdated = value.Time
			}
		case guildrole.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field guild_roles", value)
			} else if value.Valid {
				gr.guild_roles = new(uint64)
				*gr.guild_roles = uint64(value.Int64)
			}
		case guildrole.ForeignKeys[1]:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field you_tube_talent_roles", values[i])
			} else if value.Valid {
				gr.you_tube_talent_roles = new(string)
				*gr.you_tube_talent_roles = value.String
			}
		}
	}
	return nil
}

// QueryGuild queries the "guild" edge of the GuildRole entity.
func (gr *GuildRole) QueryGuild() *GuildQuery {
	return (&GuildRoleClient{config: gr.config}).QueryGuild(gr)
}

// QueryUserMemberships queries the "user_memberships" edge of the GuildRole entity.
func (gr *GuildRole) QueryUserMemberships() *UserMembershipQuery {
	return (&GuildRoleClient{config: gr.config}).QueryUserMemberships(gr)
}

// QueryTalent queries the "talent" edge of the GuildRole entity.
func (gr *GuildRole) QueryTalent() *YouTubeTalentQuery {
	return (&GuildRoleClient{config: gr.config}).QueryTalent(gr)
}

// Update returns a builder for updating this GuildRole.
// Note that you need to call GuildRole.Unwrap() before calling this method if this GuildRole
// was returned from a transaction, and the transaction was committed or rolled back.
func (gr *GuildRole) Update() *GuildRoleUpdateOne {
	return (&GuildRoleClient{config: gr.config}).UpdateOne(gr)
}

// Unwrap unwraps the GuildRole entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (gr *GuildRole) Unwrap() *GuildRole {
	tx, ok := gr.config.driver.(*txDriver)
	if !ok {
		panic("ent: GuildRole is not a transactional entity")
	}
	gr.config.driver = tx.drv
	return gr
}

// String implements the fmt.Stringer.
func (gr *GuildRole) String() string {
	var builder strings.Builder
	builder.WriteString("GuildRole(")
	builder.WriteString(fmt.Sprintf("id=%v", gr.ID))
	builder.WriteString(", name=")
	builder.WriteString(gr.Name)
	builder.WriteString(", last_updated=")
	builder.WriteString(gr.LastUpdated.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// GuildRoles is a parsable slice of GuildRole.
type GuildRoles []*GuildRole

func (gr GuildRoles) config(cfg config) {
	for _i := range gr {
		gr[_i].config = cfg
	}
}
