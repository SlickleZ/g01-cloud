function increase(quantity) {
    var old = document.getElementById("number").innerHTML;
    if (parseInt(old) + 1 >= quantity) {
        document.getElementById("number").innerHTML = parseInt(quantity);
    }else {
        document.getElementById("number").innerHTML = parseInt(old) + 1;
    }
}
    

function decrease() {
    var old = document.getElementById("number").innerHTML;
    if (parseInt(old) - 1 > 0) {
        document.getElementById("number").innerHTML = parseInt(old) - 1;
    } else {
        document.getElementById("number").innerHTML = 0;
    }
}