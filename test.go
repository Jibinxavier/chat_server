package main
import ( 
    "fmt" 
)
type Item struct {
    Id          string
    Name        string
    Price       string
}

type Cart1 struct {
    Id          string
    Items       []Item
}
type Cart2 struct {
    Id          string
    Items       []*Item
}
func main() {
    foo := cart1.Items[0]
	foo.Name := "foo" //will not change cart1
	//but in pointer case
	bar := cart2.Items[0]
	bar.Name := "bar" //will change cart2.Items[0].Name to "bar"
	fm
  

}