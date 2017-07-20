package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	BASE_LOWER_BOUND         = 2
	BASE_UPPER_BOUND         = 1024
	PSEUDO_PRIME_LOWER_BOUND = 3
	PSEUDO_PRIME_UPPER_BOUND = 10000000
	//	PSEUDO_PRIME_UPPER_BOUND = 1000000
	//	PSEUDO_PRIME_UPPER_BOUND = 100000
	PSEUDO_PRIME_RANGE_DIV = 100			//	大きい数字の時、タスクを分け合えるようにそれなりに分割する
)

type BaseNumberAndPseudoPrimeNumbers struct {
	BaseNumber         int64
	PseudoPrimeNumbers []int64
}

//	ソート用
/*
type BaseNumberAndPseudoPrimeNumbersSlice []BaseNumberAndPseudoPrimeNumbers
func (b BaseNumberAndPseudoPrimeNumbersSlice) Len() int {
	return len(b)
}
func (b BaseNumberAndPseudoPrimeNumbersSlice) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b BaseNumberAndPseudoPrimeNumbersSlice) Less(i, j int) bool {
	return b[i].BaseNumber < b[i].BaseNumber
}
*/
//	sort.Slice(b, func(i, j int) bool {
//	        return b[i].BaseNumber < b[j].BaseNumber
//	})

func main() {
	/*
		var base int64 = 2
		var exp int64 = 4
		var mod int64 = 13
		fmt.Printf("%d : %d : %d\n", base, exp, mod)
		ans := powmod(base, exp, mod)
		fmt.Printf("%d\n", ans)
	*/

	start := time.Now()

	strong_pseudoprimes := make(map[int64]*BaseNumberAndPseudoPrimeNumbers)

	mutex := new(sync.Mutex)
	wait_group := new(sync.WaitGroup)

	pseudo_prime_num := int64(PSEUDO_PRIME_UPPER_BOUND - PSEUDO_PRIME_LOWER_BOUND + 1)

	for i := 0; i < PSEUDO_PRIME_RANGE_DIV; i++ {
		wait_group.Add(1)
		go func(index int) {
			defer wait_group.Done()

			var lower int64 = int64(PSEUDO_PRIME_LOWER_BOUND) + (pseudo_prime_num/PSEUDO_PRIME_RANGE_DIV)*int64(index)
			var upper int64 = lower + (pseudo_prime_num / PSEUDO_PRIME_RANGE_DIV)
			if upper > PSEUDO_PRIME_UPPER_BOUND {
				upper = PSEUDO_PRIME_UPPER_BOUND
			}

			for probable_prime := lower; probable_prime < upper; probable_prime++ {
				if probable_prime&1 == 0 {
					continue
				}
				calc_pseudoprimes(strong_pseudoprimes, mutex, probable_prime, BASE_LOWER_BOUND, BASE_UPPER_BOUND)
			}
		}(i)
	}

	wait_group.Wait()
	end := time.Now()

	//	結果出力
	printList(strong_pseudoprimes)
	fmt.Fprintf(os.Stderr, "%.2f sec\n", float64((end.Sub(start)).Nanoseconds())/1000000000.0)

}

func printList(strong_pseudoprimes map[int64]*BaseNumberAndPseudoPrimeNumbers) {

	sorted_array := make([]*BaseNumberAndPseudoPrimeNumbers, 0, len(strong_pseudoprimes))
	for _, base_number_and_pseudo_prime_numbers := range strong_pseudoprimes {
		sorted_array = append(sorted_array, base_number_and_pseudo_prime_numbers)
	}
	sort.Slice(sorted_array, func(i, j int) bool { return sorted_array[i].BaseNumber < sorted_array[j].BaseNumber })

	for _, base_number_and_pseudo_prime_numbers := range sorted_array {

		str := ""
		sort.Slice(base_number_and_pseudo_prime_numbers.PseudoPrimeNumbers, func(i, j int) bool {
			return base_number_and_pseudo_prime_numbers.PseudoPrimeNumbers[i] < base_number_and_pseudo_prime_numbers.PseudoPrimeNumbers[j]
		})

		for _, v := range base_number_and_pseudo_prime_numbers.PseudoPrimeNumbers {
			str += fmt.Sprintf("%d, ", v)
		}
		str = strings.TrimRight(str, ", ") //右端の","を取り除く

		fmt.Printf("%4d : %4d : [ %s ]\n", base_number_and_pseudo_prime_numbers.BaseNumber, len(base_number_and_pseudo_prime_numbers.PseudoPrimeNumbers), str)

	}

}

func powmod(base int64, exp int64, mod int64) int64 {

	if base >= mod {
		base = base % mod
	}
	//	base が mod の倍数、もしくは等しい場合を除いておく
	if base == 0 {
		return 0
	}

	var answer int64 = 1
	for exp > 0 {
		if exp&1 == 1 {
			answer = (answer * base) % mod
		}
		base = (base * base) % mod
		exp >>= 1
	}
	return answer

}

func miller_rabin_primality_test(probable_prime int64, base int64) int {

	// 互いに素でないなら判定しない(片方がもう片方の倍数)
	if (base >= probable_prime && base%probable_prime == 0) || (base < probable_prime && probable_prime%base == 0) {
		return -1
	}

	// ミラーラビンテストで使用する d を求める
	// d : (probable_prime-1) = ( d*(2*r) ) を求める
	miller_rabin_d := probable_prime - 1
	for miller_rabin_d&1 == 0 {
		miller_rabin_d >>= 1
	}

	t := miller_rabin_d
	y := powmod(base, t, probable_prime)
	for t != probable_prime-1 && y != 1 && y != probable_prime-1 {
		y = (y * y) % probable_prime
		t <<= 1
	}
	if y != probable_prime-1 && t&1 == 0 {
		return 0
	}
	return 1
}

func calc_pseudoprimes(strong_pseudoprimes map[int64]*BaseNumberAndPseudoPrimeNumbers, mutex *sync.Mutex, probable_prime int64, lower_bound int64, upper_bound int64) {
	mr_test_results := make(map[int64]int, upper_bound-lower_bound+1)

	for base := lower_bound; base <= upper_bound; base++ {
		mr_test_results[base] = miller_rabin_primality_test(probable_prime, base)
	}

	answers_count := 0
	answers_sum := 0
	for _, test_result := range mr_test_results {
		if test_result >= 0 {
			answers_count++
			if test_result == 1 {
				answers_sum++
			}
		}
	}

	if answers_sum == 0 {
		//		return :composit
	} else if answers_sum == answers_count {
		//		return :prime
	} else {
		//	pseudoprime
		for base_number, test_result := range mr_test_results {
			if test_result == 1 {
				mutex.Lock()
				if _, ok := strong_pseudoprimes[base_number]; !ok {
					//	strong_pseudoprimes[base_number] = make([]int64, 0, 1000)
					strong_pseudoprimes[base_number] = new(BaseNumberAndPseudoPrimeNumbers)
					strong_pseudoprimes[base_number].BaseNumber = base_number
					strong_pseudoprimes[base_number].PseudoPrimeNumbers = make([]int64, 0, 1000)
				}
				strong_pseudoprimes[base_number].PseudoPrimeNumbers = append(strong_pseudoprimes[base_number].PseudoPrimeNumbers, probable_prime)
				mutex.Unlock()
			}

		}
	}

}
