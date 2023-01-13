package wait_demo

import (
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/clock"
	"testing"
	"time"
)

func TestPollImmediateUntil(t *testing.T) {
	stopCh := make(<-chan struct{})

	go func() {
		err := PollImmediateUntil(2*time.Second,
			func() (bool, error) {
				t.Log("被调用了")
				return false, nil
			},
			stopCh)
		if err != nil {
			t.Fatal(err)
		}
	}()

	time.Sleep(10 * time.Minute)
}

func NewTestBackoff() *testBackoff {
	return &testBackoff{Clock: clock.RealClock{}}
}

type testBackoff struct {
	Clock clock.Clock
}

func (t *testBackoff) Backoff() clock.Timer {
	return t.Clock.NewTimer(5 * time.Second)
}

func TestBackoffUntil(t *testing.T) {
	stopCh := make(<-chan struct{})
	backoffManager := NewTestBackoff()
	BackoffUntil(func() {
		t.Log("我被调用了")
	}, backoffManager, false, stopCh)
}

func TestPollImmediate(t *testing.T) {
	if err := PollImmediate(1*time.Second, 2*time.Second, func() (done bool, err error) {
		t.Logf("start")
		return false, nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestK8SBackoffUntil(t *testing.T) {
	stopCh := make(<-chan struct{})
	backoffManager := wait.NewExponentialBackoffManager(800*time.Millisecond, 30*time.Second, 2*time.Minute, 2.0, 1.0, &clock.RealClock{})
	BackoffUntil(func() {
		t.Log("我被调用了")
	}, backoffManager, false, stopCh)
}
