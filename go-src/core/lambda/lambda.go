package lambda

// 定义函数类型，模拟Lambda表达式

// Mapper 用于map操作的函数类型
type Mapper[T, R any] func(T) R

// Predicate 用于filter操作的函数类型
type Predicate[T any] func(T) bool

// Reducer 用于reduce操作的函数类型
type Reducer[T, R any] func(acc R, item T) R

// Consumer 用于forEach操作的函数类型
type Consumer[T any] func(T)

// 集合操作实现

// Map 对集合中的每个元素应用Mapper函数，返回新的集合
func Map[T, R any](collection []T, mapper Mapper[T, R]) []R {
	result := make([]R, 0, len(collection))
	for _, item := range collection {
		result = append(result, mapper(item))
	}
	return result
}

// Filter 对集合中的元素应用Predicate函数，返回满足条件的元素集合
func Filter[T any](collection []T, predicate Predicate[T]) []T {
	result := make([]T, 0)
	for _, item := range collection {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Reduce 对集合中的元素应用Reducer函数，将集合归约为单个值
func Reduce[T, R any](collection []T, reducer Reducer[T, R], initial R) R {
	acc := initial
	for _, item := range collection {
		acc = reducer(acc, item)
	}
	return acc
}

// ForEach 对集合中的每个元素应用Consumer函数
func ForEach[T any](collection []T, consumer Consumer[T]) {
	for _, item := range collection {
		consumer(item)
	}
}

// Find 查找集合中第一个满足Predicate条件的元素
func Find[T any](collection []T, predicate Predicate[T]) (T, bool) {
	for _, item := range collection {
		if predicate(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

// AnyMatch 检查集合中是否有任何元素满足Predicate条件
func AnyMatch[T any](collection []T, predicate Predicate[T]) bool {
	for _, item := range collection {
		if predicate(item) {
			return true
		}
	}
	return false
}

// AllMatch 检查集合中是否所有元素都满足Predicate条件
func AllMatch[T any](collection []T, predicate Predicate[T]) bool {
	for _, item := range collection {
		if !predicate(item) {
			return false
		}
	}
	return true
}

// NoneMatch 检查集合中是否没有元素满足Predicate条件
func NoneMatch[T any](collection []T, predicate Predicate[T]) bool {
	return !AnyMatch(collection, predicate)
}

func GroupBy[T any, K comparable](slice []T, keyFunc func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, item := range slice {
		key := keyFunc(item)
		result[key] = append(result[key], item)
	}
	return result
}
