import "add"
import "math"

//You can comment code
func main(){
   print "Log(2) = "+ln(2)
   print "Exp(2) = "+exp(2)
   print math.power(2,2)
   print math.sqrt(4,2)

   var bool = true
   if bool == true{
      print "Bool is true"
      }
   var hallo="Helo world"
   print hallo
   print add.Add(6,2)
   print add.Add(6,4)
   print "Write something:"
   var s = Read()
   {} //Throw error, but don't crash
   print 'You wrote: "' + s +'", didnt you?'
   print "Random number: " + Random(1, 5)
   print 5/2 //can use float
}
