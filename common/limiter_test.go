package common

import (
	"fmt"
	"testing"
	"time"
)

func TestLimiter(t *testing.T){
	limiter:=NewChannelLimiter(10)
	for i:=0;i<20;i++{
		if limiter.Allow(){
			go func(j int) {
				time.Sleep(time.Second)
				fmt.Println(j)
				limiter.Release()
			}(i)
		}else{
			time.Sleep(2*time.Second)
		}
	}
}
