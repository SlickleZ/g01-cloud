function checkOut() {
    number = document.getElementById("number").innerHTML
    price_view = document.getElementById("price-view").innerHTML
    total_quantity = document.getElementById("quantity-view").innerHTML
    name_view = document.getElementById("name-view").innerHTML
    detail_view = document.getElementById("detail-view").innerHTML
    
    document.getElementById("number-display").innerHTML = number
    document.getElementById("total-price-display").innerHTML = parseInt(number) * parseInt(price_view)

    document.getElementById("number_post").value = number
    document.getElementById("name_post").value = name_view
    document.getElementById("detail_post").value = detail_view
    document.getElementById("price_post").value = price_view
    document.getElementById("total_quan_post").value = total_quantity
}