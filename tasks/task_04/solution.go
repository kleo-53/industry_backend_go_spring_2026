package main

type Stats struct {
	Count int
	Sum int64
	Min int64
	Max int64
}

func Calc(nums []int64) Stats {
	if len(nums) == 0 {
		return Stats{}
	}
	res := Stats{
		Count: len(nums),
		Sum: 0,
		Min: nums[0],
		Max: nums[0],
	}
	for _, el := range nums {
		res.Sum += el
		if res.Max < el {
			res.Max = el
		}
		if res.Min > el {
			res.Min = el
		}
	}
	return res
}
