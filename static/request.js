// serRequest fetches any data from server, return json
export async function serverRequest(uri, query) {
    let response = await fetch(uri + "?" + query)
    if (!response.ok) {
        throw new Error(`HTTP error, starus: ${response.status}`)
    }
    return response.json()
}