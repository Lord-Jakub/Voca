func main() {
graphics.Init(800,800,"Ahoj")
graphics.SetFPS(60)
var x = 100
var y = 100
while isRunning() == true {
    graphics.DrawImage(x,y,"img.png")
    if graphics.KeyRight() == true{
        x = x+2
    }
    if graphics.KeyLeft() == true{
        x = x-2
    }
    if graphics.KeyUp() == true{
        y = y-2
    }
    if graphics.KeyDown() == true{
        y = y+2
    }
    graphics.Update()
    }
graphics.Close()
}