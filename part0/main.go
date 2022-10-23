package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"time"
)

func main() {
	//slice
	slice := []int{1, 2, 3, 4, 5}
	fmt.Printf("sliece: %v\n", slice)

	slice2 := make([]int, 5)
	copy(slice2, slice)
	fmt.Printf("sliece: %v, slice2: %v\n", slice, slice2)

	//delete
	i := 2
	slice = append(slice[:i], slice[i+1:]...)
	fmt.Printf("sliece: %v\n", slice)

	//reverse
	slice4 := []int{1, 2, 3, 4, 5}
	for left, right := 0, len(slice4)-1; left < right; left, right = left+1, right-1 {
		slice4[left], slice4[right] = slice4[right], slice4[left]
	}
	fmt.Printf("slice4: %v\n", slice4)

	//並行処理
	if err := doSomeThingParallel(2); err != nil {
		fmt.Printf("err: %v\n", err.Error())
	}
	fmt.Println("complete doSomeThingParallel")

	//context deadline
	// exampleDeadLine()

	//context timeout
	// exampleWithTimeout()

	//context value
	exampleWithValue()

	//pointer
	s := []struct {
		Number int
	}{
		{1}, {2}, {3}, {4}, {5},
	}
	s2 := []*struct{ Number int }{}
	for _, v := range s {
		v := v //イテレーション時の変数vは使い回しているため新しい変数を作成してあげないと全部同じポインタを指してしまう
		s2 = append(s2, &v)
	}
	for _, v := range s2 {
		fmt.Printf("%+v\n", v)
	}
}

func IsPNG(r io.Reader) (bool, error) {
	magicnum := []byte{137, 80, 78, 71}
	buf := make([]byte, len(magicnum))
	_, err := io.ReadAtLeast(r, buf, len(buf))
	if err != nil {
		return false, err
	}
	return bytes.Equal(magicnum, buf), nil
}

func doSomeThingParallel(workerNum int) error {
	ctx := context.Background()
	cancelCtx, cancel := context.WithCancel(ctx)

	//忘れずに閉じる
	defer cancel()

	errCh := make(chan error, workerNum)

	wg := sync.WaitGroup{}
	for i := 0; i < workerNum; i++ {
		i := i
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			if err := doSomeThingWithContext(cancelCtx, num); err != nil {
				cancel()
				errCh <- err
			}
		}(i)
	}

	//並行処理の終了を待つ
	wg.Wait()

	//エラーチャネルに入ったメッセージを取り出す
	close(errCh)
	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}

	//エラーがあれば最初のエラーを返す
	if len(errs) > 0 {
		return errs[0]
	}

	//正常終了
	return nil
}

func doSomeThingWithContext(ctx context.Context, num int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	//contextがまだキャンセルされていなければ、そのまま処理にすすむ
	default:
	}
	fmt.Println(num)
	return nil
}

func exampleDeadLine() {
	ctx := context.Background()
	d := time.Now().Add(10 * time.Second)
	timerCtx, cancel := context.WithDeadline(ctx, d)

	//忘れず閉じる
	defer cancel()

	//指定時刻の1日後の時刻
	nd := d.AddDate(0, 0, 1)

	select {
	case <-time.After(time.Until(nd)):
		fmt.Println("指定の時刻になりました。")
	case <-timerCtx.Done():
		fmt.Println(timerCtx.Err())
	}
}

func exampleWithTimeout() {
	ctx := context.Background()
	d := 15 * time.Second
	timerCtx, cancel := context.WithTimeout(ctx, d)

	//忘れず閉じる
	defer cancel()

	//a) 10秒経つ or contextがキャンセルされるまで待機
	select {
	case <-time.After(10 * time.Second):
		fmt.Println("10秒経ちました。")
	case <-timerCtx.Done():
		fmt.Println(timerCtx.Err())
	}

	//b) 10秒経つ or contextがキャンセルされるまで待機
	select {
	case <-time.After(10 * time.Second):
		fmt.Println("10秒経ちました。")
	case <-timerCtx.Done():
		fmt.Println(timerCtx.Err())
	}

	//10秒たってaを抜けた後15秒後にcontextがキャンセルされる
}

func exampleWithValue() {
	ctx := context.Background()
	valueCtx := context.WithValue(ctx, "key1", "value1")
	fmt.Println(valueCtx.Value("key1").(string))
}
