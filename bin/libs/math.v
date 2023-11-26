func power(x,n){
    if n == 0{
        return 1
    }
    if n == 1{
        return x
    }
    var logh = ln(x)
    var tmp = n*logh
    return exp(tmp)
}
func sqrt(x,n){
    var tmp = 1/n
    return math.power(x,tmp)
}