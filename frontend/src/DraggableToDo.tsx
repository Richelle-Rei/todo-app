import {type ToDoEntry } from './ToDoEntry'
import {useSortable} from '@dnd-kit/sortable';
import {CSS} from '@dnd-kit/utilities';

export function DraggableToDo({ id, entry, onCheck, onDelete}: {
    id: number
    entry: ToDoEntry
    onCheck: (id: number) => void
    onDelete: (id: number) => void
  }) {

  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,

  } = useSortable({id});
  
  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.7 : 1,
  };
  
  return (
    <li ref={setNodeRef} style={style} className='flex h-box p-2.5 ring bg-white ring-gray-300 rounded-sm items-center my-3 '>
      <div {...attributes} {...listeners} className='cursor-grab mr-2 p-1 s-20' onClick={(e)=>{e.stopPropagation}}>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 15 15" width="1em" height="1em">
          <path fill="currentColor" fillRule="evenodd"
          d="M5.5 4.625a1.125 1.125 0 1 0 0-2.25a1.125 1.125 0 0 0 0 2.25m4 0a1.125 1.125 0 1 0 0-2.25a1.125 1.125 0 0 0 0 2.25M10.625 7.5a1.125 1.125 0 1 1-2.25 0a1.125 1.125 0 0 1 2.25 0M5.5 8.625a1.125 1.125 0 1 0 0-2.25a1.125 1.125 0 0 0 0 2.25m5.125 2.875a1.125 1.125 0 1 1-2.25 0a1.125 1.125 0 0 1 2.25 0M5.5 12.625a1.125 1.125 0 1 0 0-2.25a1.125 1.125 0 0 0 0 2.25"
          clipRule="evenodd"
          ></path>
        </svg>
      </div>
      <input type="checkbox" checked={entry.completed} onChange={(e)=>{e.stopPropagation(); onCheck(entry.id)}} className='mx-1 size-4'/>
        
      {entry.completed? <div className='text-neutral-600 mx-2 flex-1 line-through decoration-neutral-600'>{entry.title}</div> 
        : <div className='text-neutral-600 mx-2 flex-1 '>{entry.title}</div>}

      <button onClick={(e)=>{e.stopPropagation(); onDelete(entry.id)}}>
        <svg className='size-4.5 text-red-600 mx-2' 
            xmlns="http://www.w3.org/2000/svg" fill="none" 
            viewBox="0 0 24 24" strokeWidth="2" stroke="currentColor" >
        <path  strokeLinecap="round" strokeLinejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
        </svg> 

      </button>

    </li>
              
  );
}