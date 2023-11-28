func main() {
graphics.Init(800,800,"Ahoj")
graphics.SetFPS(60)
while isRunning == true {
    graphics.DrawImage(500,500,"img.png")
    graphics.Update()
    }
graphics.Close()
}