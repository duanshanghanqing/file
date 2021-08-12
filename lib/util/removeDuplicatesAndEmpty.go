package util

func RemoveDuplicatesAndEmpty(arr []string) (ret []string){
	arrLen := len(arr)
	for i:=0; i < arrLen; i++{
		if (i > 0 && arr[i-1] == arr[i]) || len(arr[i]) == 0 {
			continue;
		}
		ret = append(ret, arr[i])
	}
	return
}
