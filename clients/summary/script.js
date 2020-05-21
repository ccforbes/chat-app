let title = document.getElementById("title")
let description = document.getElementById("description")
let images = document.getElementById("images")
let errors = document.getElementById("errors")

document.getElementById("submit-button").addEventListener("click", () => {
    let url = document.getElementById("search-bar").value
    let query = "https://api.bopboyz222.xyz/v1/summary?url="
    if (url.indexOf("http://") == -1 && url.indexOf("https://") == -1) {
        query += "https://"
    }
    title.innerHTML = ""
    description.innerHTML = ""
    images.innerHTML = ""
    errors.innerHTML = ""
    fetch(query + url)
    .then(data => {
        return data.json()
    })
    .then(json => {
        return json
    })
    .then(results => {
        title.appendChild(createNewElement(results.title, "h4"))
        description.appendChild(createNewElement(results.description, "h4"))
        if (results.images) {
            results.images.map(image => {
                images.appendChild(createNewElement(image.url, "img"))
            })
        } 
    })
    .catch(err => {
        let element = document.createElement("p")
        element.innerHTML = "The website you entered was invalid. Please try again."
        errors.appendChild(element)
    })
})

function createNewElement(resultData, tag) {
    let element = document.createElement(tag)
    if (resultData) {
        if (tag == "img") {
            element.src = resultData
        } else {
            element.innerHTML = resultData
        }
    } else {
        element.innerHTML = "No Data Found"
    }
    return element
}