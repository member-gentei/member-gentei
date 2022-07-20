// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/member-gentei/member-gentei/gentei/ent/youtubetalent"
)

// YouTubeTalent is the model entity for the YouTubeTalent schema.
type YouTubeTalent struct {
	config `json:"-"`
	// ID of the ent.
	// YouTube channel ID
	ID string `json:"id,omitempty"`
	// YouTube channel name
	ChannelName string `json:"channel_name,omitempty"`
	// URL of the talent's YouTube thumbnail
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
	// ID of a members-only video
	MembershipVideoID string `json:"membership_video_id,omitempty"`
	// Last time membership_video_id returned no results
	LastMembershipVideoIDMiss time.Time `json:"last_membership_video_id_miss,omitempty"`
	// Last time data was fetched
	LastUpdated time.Time `json:"last_updated,omitempty"`
	// When refresh/membership checks were disabled. Set to zero/nil to re-enable.
	Disabled time.Time `json:"disabled,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the YouTubeTalentQuery when eager-loading is set.
	Edges YouTubeTalentEdges `json:"edges"`
}

// YouTubeTalentEdges holds the relations/edges for other nodes in the graph.
type YouTubeTalentEdges struct {
	// Guilds holds the value of the guilds edge.
	Guilds []*Guild `json:"guilds,omitempty"`
	// Roles holds the value of the roles edge.
	Roles []*GuildRole `json:"roles,omitempty"`
	// Memberships holds the value of the memberships edge.
	Memberships []*UserMembership `json:"memberships,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// GuildsOrErr returns the Guilds value or an error if the edge
// was not loaded in eager-loading.
func (e YouTubeTalentEdges) GuildsOrErr() ([]*Guild, error) {
	if e.loadedTypes[0] {
		return e.Guilds, nil
	}
	return nil, &NotLoadedError{edge: "guilds"}
}

// RolesOrErr returns the Roles value or an error if the edge
// was not loaded in eager-loading.
func (e YouTubeTalentEdges) RolesOrErr() ([]*GuildRole, error) {
	if e.loadedTypes[1] {
		return e.Roles, nil
	}
	return nil, &NotLoadedError{edge: "roles"}
}

// MembershipsOrErr returns the Memberships value or an error if the edge
// was not loaded in eager-loading.
func (e YouTubeTalentEdges) MembershipsOrErr() ([]*UserMembership, error) {
	if e.loadedTypes[2] {
		return e.Memberships, nil
	}
	return nil, &NotLoadedError{edge: "memberships"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*YouTubeTalent) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case youtubetalent.FieldID, youtubetalent.FieldChannelName, youtubetalent.FieldThumbnailURL, youtubetalent.FieldMembershipVideoID:
			values[i] = new(sql.NullString)
		case youtubetalent.FieldLastMembershipVideoIDMiss, youtubetalent.FieldLastUpdated, youtubetalent.FieldDisabled:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type YouTubeTalent", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the YouTubeTalent fields.
func (ytt *YouTubeTalent) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case youtubetalent.FieldID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value.Valid {
				ytt.ID = value.String
			}
		case youtubetalent.FieldChannelName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field channel_name", values[i])
			} else if value.Valid {
				ytt.ChannelName = value.String
			}
		case youtubetalent.FieldThumbnailURL:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field thumbnail_url", values[i])
			} else if value.Valid {
				ytt.ThumbnailURL = value.String
			}
		case youtubetalent.FieldMembershipVideoID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field membership_video_id", values[i])
			} else if value.Valid {
				ytt.MembershipVideoID = value.String
			}
		case youtubetalent.FieldLastMembershipVideoIDMiss:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field last_membership_video_id_miss", values[i])
			} else if value.Valid {
				ytt.LastMembershipVideoIDMiss = value.Time
			}
		case youtubetalent.FieldLastUpdated:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field last_updated", values[i])
			} else if value.Valid {
				ytt.LastUpdated = value.Time
			}
		case youtubetalent.FieldDisabled:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field disabled", values[i])
			} else if value.Valid {
				ytt.Disabled = value.Time
			}
		}
	}
	return nil
}

// QueryGuilds queries the "guilds" edge of the YouTubeTalent entity.
func (ytt *YouTubeTalent) QueryGuilds() *GuildQuery {
	return (&YouTubeTalentClient{config: ytt.config}).QueryGuilds(ytt)
}

// QueryRoles queries the "roles" edge of the YouTubeTalent entity.
func (ytt *YouTubeTalent) QueryRoles() *GuildRoleQuery {
	return (&YouTubeTalentClient{config: ytt.config}).QueryRoles(ytt)
}

// QueryMemberships queries the "memberships" edge of the YouTubeTalent entity.
func (ytt *YouTubeTalent) QueryMemberships() *UserMembershipQuery {
	return (&YouTubeTalentClient{config: ytt.config}).QueryMemberships(ytt)
}

// Update returns a builder for updating this YouTubeTalent.
// Note that you need to call YouTubeTalent.Unwrap() before calling this method if this YouTubeTalent
// was returned from a transaction, and the transaction was committed or rolled back.
func (ytt *YouTubeTalent) Update() *YouTubeTalentUpdateOne {
	return (&YouTubeTalentClient{config: ytt.config}).UpdateOne(ytt)
}

// Unwrap unwraps the YouTubeTalent entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ytt *YouTubeTalent) Unwrap() *YouTubeTalent {
	_tx, ok := ytt.config.driver.(*txDriver)
	if !ok {
		panic("ent: YouTubeTalent is not a transactional entity")
	}
	ytt.config.driver = _tx.drv
	return ytt
}

// String implements the fmt.Stringer.
func (ytt *YouTubeTalent) String() string {
	var builder strings.Builder
	builder.WriteString("YouTubeTalent(")
	builder.WriteString(fmt.Sprintf("id=%v, ", ytt.ID))
	builder.WriteString("channel_name=")
	builder.WriteString(ytt.ChannelName)
	builder.WriteString(", ")
	builder.WriteString("thumbnail_url=")
	builder.WriteString(ytt.ThumbnailURL)
	builder.WriteString(", ")
	builder.WriteString("membership_video_id=")
	builder.WriteString(ytt.MembershipVideoID)
	builder.WriteString(", ")
	builder.WriteString("last_membership_video_id_miss=")
	builder.WriteString(ytt.LastMembershipVideoIDMiss.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("last_updated=")
	builder.WriteString(ytt.LastUpdated.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("disabled=")
	builder.WriteString(ytt.Disabled.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// YouTubeTalents is a parsable slice of YouTubeTalent.
type YouTubeTalents []*YouTubeTalent

func (ytt YouTubeTalents) config(cfg config) {
	for _i := range ytt {
		ytt[_i].config = cfg
	}
}
