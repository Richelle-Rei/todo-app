import { type ToDoEntry } from "./ToDoEntry"

const url = import.meta.env.VITE_BACKEND_URL + "/todos"

export async function getToDo(): Promise<ToDoEntry[]> {
    const response = await fetch(url)

    if (!response.ok) {
    const error = new Error(`Response status: ${response.status}`)
    throw error
    }
    const data = await response.json()
    console.log(data)

    return data
}

export async function createToDo(title: string, displayOrder: number): Promise<ToDoEntry> {
    const response = await fetch(url, {
        method: 'POST',
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ title, completed: false, userId: 1, displayOrder}),
      })

    if (!response.ok) {
    const error = new Error(`Response status: ${response.status}`)
    throw error
    }
    const newToDo = await response.json()
    console.log(newToDo)

    return newToDo
}

export async function updateToDo (id: number, completed: boolean): Promise<ToDoEntry> {
    const response = await fetch(`${url}/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({id, completed}),
    });
    return await response.json();
};

export async function deleteToDo (id:number) {
    const response = await fetch(`${url}/${id}`, {
      method: 'DELETE',
    });
    if(!response.ok){
        throw new Error(`Error! Response status: ${response.status}`)
    }
};

export async function updateDisplay (newList: ToDoEntry[]): Promise<ToDoEntry> {
  const response = await fetch(`${url}/order`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({newList}),
  });
  return await response.json();
};

export async function editToDo (id:number, title: string): Promise<ToDoEntry> {
  const response = await fetch(`${url}/${id}/edit`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({title}),
  });
  return await response.json();
};

