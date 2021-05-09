package common

import (
	"testing"
	"time"

	"cloud.google.com/go/firestore"
)

func TestGetPartition(t *testing.T) {
	bigSlice := make([]*firestore.DocumentSnapshot, 6000)
	for i := 0; i < len(bigSlice); i++ {
		bigSlice[i] = &firestore.DocumentSnapshot{
			CreateTime: time.Date(0, 0, 0, 0, 0, i, 0, time.UTC),
		}
	}
	partition := getPartition(bigSlice, 3, DefaultPartitionThreshold)
	if len(partition) != 2000 {
		t.Fatalf("partition of 3 did not divide evenly")
	}
	for i := 1; i < len(partition); i++ {
		if partition[i].CreateTime.Unix()-partition[i-1].CreateTime.Unix() != 3 {
			t.Fatalf("offset between items is not equal to the number of partitions")
		}
	}
	// test threshold
	smolSlice := make([]*firestore.DocumentSnapshot, 100)
	partition = getPartition(smolSlice, 3, 1000)
	if len(smolSlice) != len(partition) {
		t.Fatalf("slice smaller than threshold was partitioned")
	}
}
