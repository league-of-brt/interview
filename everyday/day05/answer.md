常规思路：
时间复杂度换取空间复杂度，构造K个元素的最小堆；
利用堆顶最小的特点，将后面大于堆顶的元素不断替换成堆顶，并且重新调整最小堆，最终堆里的元素就是TopK项。

步骤：
1）构建数组，第一个元素为堆顶
2）在堆末尾插入元素时，需要不断上浮进行调整
待调整的值比父节点小时，交换并继续上浮直到堆顶
3）替换堆顶时，需要将新替换的值不断下沉调整
将最小子节点替换到父节点，最后叶子节点存放原来的堆顶，向上调整一波重新成为最小堆

go实现：




package main


import (
   "fmt"
   "sort"
)


type TopKHeap struct {
   K    int
   data []int
}


func NewTopKHeap(k int) *TopKHeap {
   return &TopKHeap{
      K:    k,
      data: nil,
   }
}


func (h *TopKHeap) Push(item int) {
   if len(h.data) < h.K {
      h.data = append(h.data, item)
      h.upAdjust(len(h.data) - 1)
      return
   }
   if h.data[0] >= item {
      return
   }
   h.data[0] = item
   h.downAdjust()
}


func (h *TopKHeap) TopK() []int {
   return h.data
}


// 堆末尾插入元素 上浮调整
// @index 待调整位置
func (h *TopKHeap) upAdjust(index int) {
   item := h.data[index]
   for index > 0 {
      indexParent := (index - 1) >> 1
      parent := h.data[indexParent]
      if item >= parent {
         break
      }
      h.data[index] = parent
      index = indexParent
   }
   h.data[index] = item
}


// 替换堆顶 不断下沉调整
func (h *TopKHeap) downAdjust() {
   index := 0
   item := h.data[0]


   indexChild := index*2 + 1
   for indexChild < h.K {
      indexRightChild := indexChild + 1
      if indexRightChild < h.K && h.data[indexChild] >= h.data[indexRightChild] {
         indexChild = indexRightChild
      }
      h.data[index] = h.data[indexChild]
      index = indexChild
      indexChild = 2*index + 1
   }
   h.data[index] = item
   h.upAdjust(index)
}


func main() {
   data := []int{1, 8, 30, 2, 2, 3, 10, 0, 1, 5, 6, 15, 20, 1, 2}
   top := NewTopKHeap(5)
   for _, n := range data {
      top.Push(n)
   }
   fmt.Println(top.TopK())
   sort.Ints(data)
   fmt.Println(data)
}
