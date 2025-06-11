import { useState, useEffect } from 'react'
// import reactLogo from './assets/react.svg'
// import viteLogo from '/vite.svg'

type ToDoEntry = {
  description: string
  completed: boolean
}

function App() {
  const [inputText, setInputText] = useState('')
  const [toDoList, setToDoList] = useState<ToDoEntry[]>([])
  const [loading, setLoading] = useState(false)
  const count =  toDoList.filter((entry)=> !entry.completed).length

  async function getData() {
    const url = "https://jsonplaceholder.typicode.com/todos?_limit=5"
    try {
      setLoading(true)
      const response = await fetch(url)
      if (!response.ok) {
        const error = new Error(`Response status: ${response.status}`)
        throw error
      }
      const json = await response.json()
    
      const initialData = json.map((entry:any) => ({description: entry.title, completed: entry.completed}))
      setToDoList(initialData)
      console.log(initialData)

    } catch (error : any) {
      console.error(error.message)
    } finally{
      setLoading(false)
    }
  }

  useEffect(() => {
    getData()
    // console.log('on mount')
  }, [])

  function AddEntry(){
    if (!inputText.trim()){
      return
    }
    setToDoList([...toDoList, {completed: false, description: inputText.trim()}])
    setInputText('')
  }

  function DeleteEntry(index: number){
    if(!toDoList[index].completed){
    }
    const newList = toDoList.filter((_, i) => i !== index)
    setToDoList(newList)
    // console.log(toDoList)
  }

  function checkEntry(index: number){
    const newList = [...toDoList]
    newList[index].completed = !newList[index].completed
    setToDoList(newList)
  }

  return (
    <>
      <div className='flex h-screen bg-slate-50'>
        <div className='w-md h-min shadow-md rounded-lg p-6 ring ring-gray-300 m-auto bg-white'>
        
          <div className='flex place-content-center'>
            <div className='my-6 text-2xl font-semibold'> Simple To Do App</div>
          </div>

          <div className='flex place-content-center'>
              <input id='ToDoInput' value={inputText}
              onChange={(text) => setInputText(text.target.value)}
              placeholder='Add a new to do...'
              onKeyUp={(e) => e.key == 'Enter' && AddEntry()}
              className='flex-1 w-xs p-2 ring rounded-sm ring-gray-300 text-neutral-600'></input>
              
              <button id='AddEntryButton' onClick={AddEntry} 
              className='flex-none py-2 px-4 rounded-sm ml-3 bg-blue-500 text-white hover:bg-blue-700 '>+ Add</button>
          </div>

          <div className='text-neutral-600 my-3.5'> {count>0? count + " active tasks":""}</div>

          <div className='flex place-content-center'>
            {!loading && (toDoList.length>0?null:<div className='text-neutral-500 my-10'>No To Dos yet. Add one above!</div>)}

            
            {loading && <div className='my-10'><svg className="size-6 animate-spin text-slate-500 mx-auto;" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth="2" stroke="currentColor" >
            <path strokeLinecap="round" strokeLinejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182m0-4.991v4.99" />
            </svg></div>}

          </div> 

          <ul className='mt-4'>
            {toDoList.map((entry, index) => (
              <li key={index} className='flex p-2.5 ring ring-gray-300 rounded-sm items-center my-3 '>
                <input type="checkbox" checked={entry.completed} onChange={()=>checkEntry(index)} className='mx-1 size-4'/>
                {entry.completed? <div className='text-neutral-600 mx-2 flex-1 line-through decoration-neutral-600'>{entry.description}</div> 
                : <div className='text-neutral-600 mx-2 flex-1 '>{entry.description}</div>}

                <button onClick={()=>DeleteEntry(index)}>
                  <svg className='size-4.5 text-red-600 mx-2' 
                  xmlns="http://www.w3.org/2000/svg" fill="none" 
                  viewBox="0 0 24 24" strokeWidth="2" stroke="currentColor" >
                  <path  strokeLinecap="round" strokeLinejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
                  </svg> 

                </button>

                </li>
              ))}
          </ul>

        </div>
      </div>
      
    </>
  )
}
export default App


