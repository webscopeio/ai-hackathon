package analyzer

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestGetContent(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	urls := []string{
		"jakub.kr",
		"jakub.kr/components/cards",
	}

	res, err := GetContent(ctx, urls)
	if err != nil {
		t.Fatalf("err=%v", err)
	}

	fmt.Printf("res=%v", res)
}
