function give_rating(rate){
    document.getElementById("rate_value").value=rate;
    for(var i=1;i<=5;i++){
        document.getElementById("rating-"+i).src="https://storage.googleapis.com/cloud-project-g01/shiba-rating-dark.png";
    }
    for(var i=1;i<=rate;i++){
        document.getElementById("rating-"+i).src="https://storage.googleapis.com/cloud-project-g01/shiba-rating-shine.png";
    }
}