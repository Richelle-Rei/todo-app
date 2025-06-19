import { useState, useEffect } from 'react'
import { type ToDoEntry } from './ToDoEntry'
import {  getToDo, createToDo, updateToDo, deleteToDo, updateDisplay, editToDo } from './fetchcalls'
import {  DndContext, type DragEndEvent } from '@dnd-kit/core';
import {  DraggableToDo } from './DraggableToDo';
import {  SortableContext, arrayMove  } from '@dnd-kit/sortable';
import { restrictToVerticalAxis } from '@dnd-kit/modifiers';


function App() {
  const [inputText, setInputText] = useState('')
  const [toDoList, setToDoList] = useState<ToDoEntry[]>([])
  const [loading, setLoading] = useState(false)
  const count =  toDoList.filter((entry)=> !entry.completed).length

  async function handleDragEnd(event : DragEndEvent) {
    const {active, over} = event;

    if (active.id !== over?.id) {
      setToDoList((items) => {
        const oldIndex = items.findIndex(item => item.id == active.id)
        const newIndex = items.findIndex(item => item.id == over?.id)

        const newToDoList = arrayMove(items, oldIndex, newIndex)

        try{
          updateDisplay(newToDoList.map((entry, index) => ({...entry, displayOrder:index+1})))
        } catch(error : any){
          console.error(error.message)
        }finally{
        return newToDoList.map((entry, index) => ({...entry, displayOrder:index+1}))
        }
      })
    }
  }
  
  async function getData() {
    try{
      setLoading(true)
      const initialData = await getToDo()
      setToDoList(initialData)
    }catch(error : any){
      console.error(error.message)
    }finally{
      setLoading(false)
    }
  }

  useEffect(() => {
    getData()
  }, [])

  async function AddEntry(){
    if (!inputText.trim()){
      return
    }

    try{
      setLoading(true)
      await createToDo(inputText, toDoList.length+1)
      await getData()
      setInputText('')
    }catch(error : any){
      console.error(error.message)
    }finally{
      setLoading(false)

  }}

  async function deleteEntry(index: number){
    try{
      await deleteToDo(index)
      await getData()
    }catch(error : any){
      console.error(error.message)
    }
  }

  async function checkEntry(index: number){
    try{
      const todo = toDoList.find(entry => entry.id ===index)
      if(!todo){
        return
      }
      await updateToDo(index, !todo.completed)
      await getData()
    }catch(error : any){
      console.error(error.message)
    }
  }

  async function editEntry(index: number, newtitle: string) {
    try{
      const todo = toDoList.find(entry => entry.id == index)
      if(!todo || !newtitle.trim() || newtitle == todo.title){
        return
      }
      await editToDo(index, newtitle)
      await getData()
    }catch(error : any){
      console.error(error.message)
    }
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
              onKeyUp={(e) => e.key == 'Enter' && AddEntry() }
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

          <DndContext modifiers={[restrictToVerticalAxis]} onDragEnd={handleDragEnd}>
          <SortableContext items={toDoList}>
          <ul className='mt-4'>
            {toDoList.map((entry, _) => (<DraggableToDo key={entry.id} id ={entry.id} entry = {entry} onCheck={checkEntry} onDelete = {deleteEntry} onEdit={editEntry}/>
              ))}
          </ul>
          </SortableContext>
          </DndContext>
          

        </div>

        
      </div>
      
    </>
  )
}

export default App


