package algo

func BinarySearch(arr []int, val, low, high, mid int) (int, bool) {
	if low > high {
		return -1, false
	} else {
		if val > arr[mid] {
			low = mid + 1
		} else if val < arr[mid] {
			high = mid - 1
		} else {
			return mid, true
		}
		mid = (low + high) / 2
	}
	return BinarySearch(arr, val, low, high, mid)
}


func LinearSearch(arr []int, val int) (int, bool) {
	lenArr := len(arr)
	for i:=0; i<lenArr; i++  {
		if val == arr[i] {
			return i, true
		}
	}
	return -1, false
}