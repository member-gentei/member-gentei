// Code generated by entc, DO NOT EDIT.

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
	// ChannelName holds the value of the "channel_name" field.
	// YouTube channel name
	ChannelName string `json:"channel_name,omitempty"`
	// ThumbnailURL holds the value of the "thumbnail_url" field.
	// URL of the talent's YouTube thumbnail
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
	// LastUpdated holds the value of the "last_updated" field.
	// Last time data was fetched
	LastUpdated time.Time `json:"last_updated,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the YouTubeTalentQuery when eager-loading is set.
	Edges                    YouTubeTalentEdges `json:"edges"`
	user_youtube_memberships *uint64
}

// YouTubeTalentEdges holds the relations/edges for other nodes in the graph.
type YouTubeTalentEdges struct {
	// Guilds holds the value of the guilds edge.
	Guilds []*Guild `json:"guilds,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// GuildsOrErr returns the Guilds value or an error if the edge
// was not loaded in eager-loading.
func (e YouTubeTalentEdges) GuildsOrErr() ([]*Guild, error) {
	if e.loadedTypes[0] {
		return e.Guilds, nil
	}
	return nil, &NotLoadedError{edge: "guilds"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*YouTubeTalent) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case youtubetalent.FieldID, youtubetalent.FieldChannelName, youtubetalent.FieldThumbnailURL:
			values[i] = new(sql.NullString)
		case youtubetalent.FieldLastUpdated:
			values[i] = new(sql.NullTime)
		case youtubetalent.ForeignKeys[0]: // user_youtube_memberships
			values[i] = new(sql.NullInt64)
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
		case youtubetalent.FieldLastUpdated:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field last_updated", values[i])
			} else if value.Valid {
				ytt.LastUpdated = value.Time
			}
		case youtubetalent.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field user_youtube_memberships", value)
			} else if value.Valid {
				ytt.user_youtube_memberships = new(uint64)
				*ytt.user_youtube_memberships = uint64(value.Int64)
			}
		}
	}
	return nil
}

// QueryGuilds queries the "guilds" edge of the YouTubeTalent entity.
func (ytt *YouTubeTalent) QueryGuilds() *GuildQuery {
	return (&YouTubeTalentClient{config: ytt.config}).QueryGuilds(ytt)
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
	tx, ok := ytt.config.driver.(*txDriver)
	if !ok {
		panic("ent: YouTubeTalent is not a transactional entity")
	}
	ytt.config.driver = tx.drv
	return ytt
}

// String implements the fmt.Stringer.
func (ytt *YouTubeTalent) String() string {
	var builder strings.Builder
	builder.WriteString("YouTubeTalent(")
	builder.WriteString(fmt.Sprintf("id=%v", ytt.ID))
	builder.WriteString(", channel_name=")
	builder.WriteString(ytt.ChannelName)
	builder.WriteString(", thumbnail_url=")
	builder.WriteString(ytt.ThumbnailURL)
	builder.WriteString(", last_updated=")
	builder.WriteString(ytt.LastUpdated.Format(time.ANSIC))
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