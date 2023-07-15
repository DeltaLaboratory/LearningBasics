import {goto} from "$app/navigation";

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

export async function list_articles() {
    let response = await fetch("http://localhost/articles", {
        method: 'GET',
    })
    if (response.status > 399) {
        throw new Error(await response.text());
    }
    return response.json();
}

export async function get_article(id: number) {
    let response = await fetch(`http://localhost/articles/${id}`, {
        method: 'GET',
    })
    if (response.status > 399) {
        throw new Error(await response.text());
    }
    return response.json();
}

export async function create_article(title: string, content: string) {
    let response = await fetch("http://localhost/articles/", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': localStorage.getItem("agnosco") as string
        },
        body: JSON.stringify({
            title: title,
            content: content,
        })
    })
    if (response.status > 399) {
        throw new Error(await response.text());
    }
    return response.json();
}

export function parseJwt (token: string): Object {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(window.atob(base64).split('').map(function (c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));

    return JSON.parse(jsonPayload);
}

export function sessionState(): boolean {
    let accessToken = localStorage.getItem("agnosco")
    if (accessToken == null) {
        return false
    }
    try {
        let claims: any = parseJwt(accessToken)
        console.log(JSON.stringify(claims))
        return (claims['exp'] * 1000) >= Date.now();
    } catch (e) {
        return false
    }
}

export function sessionAud(): string {
    let accessToken = localStorage.getItem("agnosco")
    if (accessToken == null) {
        return ""
    }
    try {
        let claims: any = parseJwt(accessToken)
        return claims['aud'];
    } catch (e) {
        return ""
    }
}