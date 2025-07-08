package membership

import (
	"context"
	"math/rand"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/member-gentei/member-gentei/gentei/ent/enttest"
)

func TestCreateMissingUserMemberships(t *testing.T) {
	const (
		channelID = "UCuwu"
	)
	var (
		ctx         = context.Background()
		db          = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
		userID      = uint64(rand.Int63())
		guildID     = uint64(rand.Int63())
		guildRoleID = uint64(rand.Int63())
		checkResult = CheckResult{
			ChannelID: channelID,
			Time:      time.Now(),
		}
	)
	defer db.Close()
	db.User.Create().
		SetID(userID).
		SetFullName("testUser#1234").
		SetAvatarHash("aaaa").
		ExecX(ctx)
	db.YouTubeTalent.Create().
		SetID(channelID).
		SetChannelName("test channel name").
		SetThumbnailURL("https://picsum.photos/200").
		ExecX(ctx)
	// assert that it doesn't create a UserMembership object unless there's a role
	err := createMissingUserMemberships(ctx, db, userID, checkResult)
	if err != nil {
		t.Fatal(err)
	}
	if count := db.UserMembership.Query().CountX(ctx); count != 0 {
		t.Fatal("UserMembership created without a corresponding role")
	}
	// assert that it creates one if there is
	db.Guild.Create().
		SetID(guildID).
		SetName("test guild").
		AddMemberIDs(userID).
		ExecX(ctx)
	db.GuildRole.Create().
		SetID(guildRoleID).
		SetName("namae").
		SetGuildID(guildID).
		SetTalentID(channelID).
		ExecX(ctx)
	err = createMissingUserMemberships(ctx, db, userID, checkResult)
	if err != nil {
		t.Fatal(err)
	}
	if count := db.UserMembership.Query().CountX(ctx); count != 1 {
		t.Fatalf("non-1 UserMembership created without a corresponding role: %d", count)
	}
	// do it again!
	err = createMissingUserMemberships(ctx, db, userID, checkResult)
	if err != nil {
		t.Fatal(err)
	}
	if count := db.UserMembership.Query().CountX(ctx); count > 1 {
		t.Fatalf("UserMembership created when it shouldn't have: %d", count)
	}
}
