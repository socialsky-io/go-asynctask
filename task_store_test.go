package taskstore_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	taskstore "github.com/haitch/taskStore"
	"github.com/stretchr/testify/assert"
)

func countingTask(ctx context.Context) (interface{}, error) {
	result := 0
	for i := 0; i < 10; i++ {
		select {
		case <-time.After(200 * time.Millisecond):
			fmt.Printf("  working %d\n", i)
			result = i
		case <-ctx.Done():
			fmt.Println("work canceled")
			return result, nil
		}
	}
	return result, nil
}

func TestEasyCase(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	t1 := taskstore.StartTask(ctx, countingTask)

	assert.Equal(t, taskstore.StateRunning, t1.State, "Task should queued to Running")

	time.Sleep(time.Second * 3)

	assert.Equal(t, taskstore.StateCompleted, t1.State, "Task should complete by now")
	assert.NotNil(t, t1.Result)
	result := t1.Result.(int)
	assert.Equal(t, result, 9)
}

func TestCancelFunc(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	t1 := taskstore.StartTask(ctx, countingTask)

	assert.Equal(t, taskstore.StateRunning, t1.State, "Task should queued to Running")

	time.Sleep(time.Second * 1)
	t1.Cancel()

	assert.Equal(t, taskstore.StateCompleted, t1.State, "Task should complete by now")
	assert.NotNil(t, t1.Result)
	result := t1.Result.(int)
	assert.Less(t, result, 9)

	// cancel a task shouldn't cancel it's parent context.
	select {
	case <-time.After(2 * time.Millisecond):
		fmt.Println("parent ctx still running")
	case <-ctx.Done():
		fmt.Println("parent ctx got canceled")
	}
}

func TestWaitTask(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	t1 := taskstore.StartTask(ctx, countingTask)

	assert.Equal(t, taskstore.StateRunning, t1.State, "Task should queued to Running")

	t1.Wait()

	assert.Equal(t, taskstore.StateCompleted, t1.State, "Task should complete by now")
	assert.NotNil(t, t1.Result)
	result := t1.Result.(int)
	assert.Equal(t, result, 9)
}

func TestDummy(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctx.Done()
}
