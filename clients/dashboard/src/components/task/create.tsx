'use client';

import { useState } from 'react';
import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Editor } from "@monaco-editor/react"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"

import { toast } from 'react-hot-toast'; // Assuming you're using react-hot-toast for notifications
import {
  createConnectTransport,
} from '@connectrpc/connect-web'
import {
  createPromiseClient,
} from '@connectrpc/connect'


import { TaskManagementService } from "@buf/evalsocket_cloud.connectrpc_es/cloud/v1/cloud_connect"
import { createTaskManagementClient } from '@/components/task/util';


// Custom hook for form state management
function useTaskForm(initialTask: any) {
  const [name, setName] = useState(initialTask?.name || '');
  const [description, setDescription] = useState(initialTask?.description || '');
  const [type, setType] = useState(initialTask?.type || 'run_query');
  const [jsonInput, setJsonInput] = useState(JSON.stringify(initialTask?.payload?.parameters || {}, null, 2));

  return { name, setName, description, setDescription, type, setType, jsonInput, setJsonInput };
}

/**
 * CreateTask component for creating a new task.
 * This component renders a modal sheet with a form to create a new task.
 */
export function CreateTask({ isOpen, setOpen, initialTask, onCreateTask }: { isOpen: boolean, setOpen: (open: boolean) => void, initialTask: any, onCreateTask: (task: any) => void }) {
  const { name, setName, description, setDescription, type, setType, jsonInput, setJsonInput } = useTaskForm(initialTask);

  const handleSubmit = async () => {
    try {
      const createdTask = await createTaskManagementClient().createTask({
        name: name,
        type: type,
        description: description,
        payload: {
          parameters: JSON.parse(jsonInput)
        },
      });

      console.log('Server response:', createdTask);

      if (createdTask && createdTask.id) {
        onCreateTask(createdTask);
        setOpen(false);
        toast.success('Task created successfully');
      } else {
        throw new Error('Task creation response is incomplete');
      }
    } catch (error) {
      console.error('Error creating task:', error);
      let errorMessage = 'Failed to create task';
      if (error instanceof Error) {
        errorMessage += `: ${error.message}`;
      }
      toast.error(errorMessage);
    }
  };

  return (
    <Sheet open={isOpen} onOpenChange={setOpen}>
      <SheetContent side="right" className={'min-w-[80vw]'}>
        <SheetHeader>
          <SheetTitle>Create New Task</SheetTitle>
          <SheetDescription>
            Enter task details and click save to create a new task.
          </SheetDescription>
        </SheetHeader>
        <div className="grid gap-4 py-4">
          <div className="grid gap-2">
            <Label htmlFor="name">Name</Label>
            <Input
              id="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              type="text"
              placeholder="Task Name"
              required
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="description">Description</Label>
            <Input
              id="description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              type="text"
              placeholder="Task Description"
              required
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="type">Type</Label>
            <Select value={type} onValueChange={setType}>
              <SelectTrigger>
                <SelectValue placeholder="Select task type" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="run_query">Run Query</SelectItem>
                <SelectItem value="send_email">Send Email</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div className="grid gap-2">
            <Label htmlFor="jsonInput">JSON Input</Label>
            <Editor
              height="300px"
              language="json"
              value={jsonInput}
              onChange={(value) => setJsonInput(value || '{}')}
              options={{ minimap: { enabled: false } }}
            />
          </div>
        </div>
        <SheetFooter>
          <Button type="submit" onClick={handleSubmit}>Create Task</Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}