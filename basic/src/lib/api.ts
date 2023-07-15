export async function register(username: string, password: string) {
    let response = await fetch("http://localhost/account/register", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            username: username,
            password: password,
        })
    })
    if (response.status > 399) {
        throw new Error(await response.text());
    }
    return response.json();
}

export async function login(username: string, password: string) {
    let response = await fetch("http://localhost/account/login", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            username: username,
            password: password,
        })
    })
    if (response.status > 399) {
        throw new Error(await response.text());
    }
    let token = await response.json();
    localStorage.setItem('agnosco', token["token"]);
    return token;
}
