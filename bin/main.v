func main() {
    graphics.Init(1000,1000,"Pong")
    graphics.SetFPS(30)
    var x = 100
    var y = 100
    var cjmp = true
    var gravity = 4
    var isjumping = false

    while isRunning() == true {
        graphics.DrawImage(x,y,"player.png")
        if graphics.KeyRight() == true{
            x = x+2
        }
        if graphics.KeyLeft() == true{
            x = x-2
        }
        if graphics.KeyUp() == true{
            if cjmp == true{
                cjmp = false
                
                isjumping = true
                gravity = -15
            }
            
        }
        if isjumping == true{
            gravity = gravity+1
            }
        if gravity > 0{
            isjumping = false
            }
        
        print gravity
        y = y+gravity
    
        if y > 600{
            if x <700{
            cjmp = true
            isjumping = false
            gravity = 0
            }
        }
    
        if y < 600{
            
            if isjumping == false{
                gravity = 4
            
            }
        }
        if x > 700{
            
            if isjumping == false{
                gravity = 4
            
             }
        }
        print isjumping
    
    
        graphics.Update()
        }
    graphics.Close()
}