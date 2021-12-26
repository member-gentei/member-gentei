// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// GuildsColumns holds the columns for the "guilds" table.
	GuildsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUint64, Increment: true},
		{Name: "name", Type: field.TypeString},
		{Name: "icon_hash", Type: field.TypeString, Nullable: true},
		{Name: "audit_channel", Type: field.TypeUint64, Unique: true, Nullable: true},
		{Name: "language", Type: field.TypeEnum, Enums: []string{"en-US"}, Default: "en-US"},
		{Name: "admin_snowflakes", Type: field.TypeJSON},
		{Name: "moderator_snowflakes", Type: field.TypeJSON, Nullable: true},
		{Name: "settings", Type: field.TypeJSON, Nullable: true},
	}
	// GuildsTable holds the schema information for the "guilds" table.
	GuildsTable = &schema.Table{
		Name:       "guilds",
		Columns:    GuildsColumns,
		PrimaryKey: []*schema.Column{GuildsColumns[0]},
	}
	// GuildRolesColumns holds the columns for the "guild_roles" table.
	GuildRolesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUint64, Increment: true},
		{Name: "name", Type: field.TypeString},
		{Name: "last_updated", Type: field.TypeTime},
		{Name: "guild_roles", Type: field.TypeUint64, Nullable: true},
	}
	// GuildRolesTable holds the schema information for the "guild_roles" table.
	GuildRolesTable = &schema.Table{
		Name:       "guild_roles",
		Columns:    GuildRolesColumns,
		PrimaryKey: []*schema.Column{GuildRolesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "guild_roles_guilds_roles",
				Columns:    []*schema.Column{GuildRolesColumns[3]},
				RefColumns: []*schema.Column{GuildsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUint64, Increment: true},
		{Name: "full_name", Type: field.TypeString},
		{Name: "avatar_hash", Type: field.TypeString},
		{Name: "last_check", Type: field.TypeTime},
		{Name: "youtube_id", Type: field.TypeString, Unique: true, Nullable: true},
		{Name: "youtube_token", Type: field.TypeJSON, Nullable: true},
		{Name: "membership_metadata", Type: field.TypeJSON, Nullable: true},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
	}
	// YouTubeTalentsColumns holds the columns for the "you_tube_talents" table.
	YouTubeTalentsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "channel_name", Type: field.TypeString},
		{Name: "thumbnail_url", Type: field.TypeString},
		{Name: "last_updated", Type: field.TypeTime},
		{Name: "user_youtube_memberships", Type: field.TypeUint64, Nullable: true},
	}
	// YouTubeTalentsTable holds the schema information for the "you_tube_talents" table.
	YouTubeTalentsTable = &schema.Table{
		Name:       "you_tube_talents",
		Columns:    YouTubeTalentsColumns,
		PrimaryKey: []*schema.Column{YouTubeTalentsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "you_tube_talents_users_youtube_memberships",
				Columns:    []*schema.Column{YouTubeTalentsColumns[4]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// GuildMembersColumns holds the columns for the "guild_members" table.
	GuildMembersColumns = []*schema.Column{
		{Name: "guild_id", Type: field.TypeUint64},
		{Name: "user_id", Type: field.TypeUint64},
	}
	// GuildMembersTable holds the schema information for the "guild_members" table.
	GuildMembersTable = &schema.Table{
		Name:       "guild_members",
		Columns:    GuildMembersColumns,
		PrimaryKey: []*schema.Column{GuildMembersColumns[0], GuildMembersColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "guild_members_guild_id",
				Columns:    []*schema.Column{GuildMembersColumns[0]},
				RefColumns: []*schema.Column{GuildsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "guild_members_user_id",
				Columns:    []*schema.Column{GuildMembersColumns[1]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// GuildAdminsColumns holds the columns for the "guild_admins" table.
	GuildAdminsColumns = []*schema.Column{
		{Name: "guild_id", Type: field.TypeUint64},
		{Name: "user_id", Type: field.TypeUint64},
	}
	// GuildAdminsTable holds the schema information for the "guild_admins" table.
	GuildAdminsTable = &schema.Table{
		Name:       "guild_admins",
		Columns:    GuildAdminsColumns,
		PrimaryKey: []*schema.Column{GuildAdminsColumns[0], GuildAdminsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "guild_admins_guild_id",
				Columns:    []*schema.Column{GuildAdminsColumns[0]},
				RefColumns: []*schema.Column{GuildsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "guild_admins_user_id",
				Columns:    []*schema.Column{GuildAdminsColumns[1]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// UserRolesColumns holds the columns for the "user_roles" table.
	UserRolesColumns = []*schema.Column{
		{Name: "user_id", Type: field.TypeUint64},
		{Name: "guild_role_id", Type: field.TypeUint64},
	}
	// UserRolesTable holds the schema information for the "user_roles" table.
	UserRolesTable = &schema.Table{
		Name:       "user_roles",
		Columns:    UserRolesColumns,
		PrimaryKey: []*schema.Column{UserRolesColumns[0], UserRolesColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "user_roles_user_id",
				Columns:    []*schema.Column{UserRolesColumns[0]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "user_roles_guild_role_id",
				Columns:    []*schema.Column{UserRolesColumns[1]},
				RefColumns: []*schema.Column{GuildRolesColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// YouTubeTalentGuildsColumns holds the columns for the "you_tube_talent_guilds" table.
	YouTubeTalentGuildsColumns = []*schema.Column{
		{Name: "you_tube_talent_id", Type: field.TypeString},
		{Name: "guild_id", Type: field.TypeUint64},
	}
	// YouTubeTalentGuildsTable holds the schema information for the "you_tube_talent_guilds" table.
	YouTubeTalentGuildsTable = &schema.Table{
		Name:       "you_tube_talent_guilds",
		Columns:    YouTubeTalentGuildsColumns,
		PrimaryKey: []*schema.Column{YouTubeTalentGuildsColumns[0], YouTubeTalentGuildsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "you_tube_talent_guilds_you_tube_talent_id",
				Columns:    []*schema.Column{YouTubeTalentGuildsColumns[0]},
				RefColumns: []*schema.Column{YouTubeTalentsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "you_tube_talent_guilds_guild_id",
				Columns:    []*schema.Column{YouTubeTalentGuildsColumns[1]},
				RefColumns: []*schema.Column{GuildsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		GuildsTable,
		GuildRolesTable,
		UsersTable,
		YouTubeTalentsTable,
		GuildMembersTable,
		GuildAdminsTable,
		UserRolesTable,
		YouTubeTalentGuildsTable,
	}
)

func init() {
	GuildRolesTable.ForeignKeys[0].RefTable = GuildsTable
	YouTubeTalentsTable.ForeignKeys[0].RefTable = UsersTable
	GuildMembersTable.ForeignKeys[0].RefTable = GuildsTable
	GuildMembersTable.ForeignKeys[1].RefTable = UsersTable
	GuildAdminsTable.ForeignKeys[0].RefTable = GuildsTable
	GuildAdminsTable.ForeignKeys[1].RefTable = UsersTable
	UserRolesTable.ForeignKeys[0].RefTable = UsersTable
	UserRolesTable.ForeignKeys[1].RefTable = GuildRolesTable
	YouTubeTalentGuildsTable.ForeignKeys[0].RefTable = YouTubeTalentsTable
	YouTubeTalentGuildsTable.ForeignKeys[1].RefTable = GuildsTable
}
