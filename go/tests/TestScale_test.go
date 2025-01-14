//go:build scale

package tests

import (
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/tests"
	"github.com/saichler/shared/go/types"
	"testing"
	"time"
)

func scaleTest(size, exp int, timeout int64, t *testing.T) bool {
	start := time.Now().Unix()
	for i := 0; i < size; i++ {
		pb := &tests.TestProto{}
		pb.Int32 = int32(i)
		err := eg2.Do(types.Action_POST, eg3.Config().Local_Uuid, pb)
		if err != nil {
			interfaces.Fail(t, err)
			return false
		}
	}

	eg3 := tsps["eg3"]

	now := time.Now().Unix()
	for eg3.PostNumber < exp {
		if time.Now().Unix()-timeout >= now {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
	end := time.Now().Unix()
	interfaces.Info("Scale test for ", size, " took ", (end - start), " seconds")
	if eg3.PostNumber != exp {
		interfaces.Fail(t, "eg3", " Post count does not equal to ", exp, ":", eg3.PostNumber)
		return false
	}
	return true
}

func TestScale(t *testing.T) {
	interfaces.Logger().SetLogLevel(interfaces.Info_Level)
	exp := 1000
	ok := scaleTest(1000, exp, 2, t)
	if !ok {
		return
	}
	exp += 10000
	ok = scaleTest(10000, exp, 2, t)
	if !ok {
		return
	}
	exp += 100000
	ok = scaleTest(100000, exp, 5, t)
	if !ok {
		return
	}
	exp += 1000000
	scaleTest(1000000, exp, 5, t)
}