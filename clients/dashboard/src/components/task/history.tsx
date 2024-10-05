'use client';


import toast from "react-hot-toast"
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

import { formatDistanceToNow } from 'date-fns'
import { useEffect, useState } from "react";  // Import useState
import { createTaskManagementClient, getStatusColor, getStatusString, type TaskHistory, type TaskInterface } from "@/components/task/util";



/**
 * TaskHistoryList component to display the details and history of a selected task.
 * @param {Object} props - Component props
 * @param {TaskInterface} props.task - The selected task
 * @param {boolean} props.isOpen - Modal open state
 * @param {Function} props.setOpen - Function to set modal open state
 */
export function TaskHistoryList({ task, isOpen, setOpen }: { task: TaskInterface; isOpen: boolean; setOpen: (open: boolean) => void; }) {
  const [taskHistory, setTaskHistory] = useState<TaskHistory[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchTaskHistory = async () => {
      setIsLoading(true);
      try {
        const response = await createTaskManagementClient().getTaskHistory({ id: task.id });

        const history = response.history.map((item: any) => ({
          id: item.id,
          details: item.details,
          status: getStatusString(item.status),
          created_at: new Date(item.createdAt).toISOString(),
        }));

        setTaskHistory(history);
      } catch (error) {
        console.error('Error fetching task history:', error);
        toast.error("Failed to fetch task history. Please try again.");
      } finally {
        setIsLoading(false);
      }
    };

    if (isOpen) {
      fetchTaskHistory();
    }
  }, [task.id, isOpen]);

  return (
    <Sheet open={isOpen} onOpenChange={setOpen}>
      <SheetContent className="sm:max-w-[600px] w-full">
        <SheetHeader>
          <SheetTitle>Task Details and History</SheetTitle>
          <SheetDescription>
            View the task details and its history of status changes.
          </SheetDescription>
        </SheetHeader>

        <div className="mt-6 space-y-4 overflow-y-auto max-h-[calc(100vh-200px)]">
          {/* Task Details */}
          <div className="bg-gray-50 dark:bg-gray-800 rounded-lg shadow-md p-6 border border-gray-200 dark:border-gray-700">
            <h3 className="font-semibold text-xl text-gray-800 dark:text-gray-200 mb-4">Task Information</h3>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-sm font-medium text-gray-500 dark:text-gray-400">ID</p>
                <p className="text-gray-900 dark:text-gray-100">{task.id}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-500 dark:text-gray-400">Name</p>
                <p className="text-gray-900 dark:text-gray-100">{task.name}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-500 dark:text-gray-400">Status</p>
                <p className={`text-gray-900 dark:text-gray-100 ${getStatusColor(getStatusString(task.status))}`}>
                  {getStatusString(task.status)}
                </p>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-500 dark:text-gray-400">Retries</p>
                <p className="text-gray-900 dark:text-gray-100">{task.retries}</p>
              </div>
            </div>
          </div>

          {isLoading ? (
            <p className="text-center">Loading task history...</p>
          ) : taskHistory.length > 0 ? (
            taskHistory.map((task: TaskHistory) => (
              <div key={task.id} className="bg-gray-100 dark:bg-gray-700 rounded-lg shadow-md p-6 transition-all duration-300 hover:shadow-lg border border-gray-200 dark:border-gray-600">
                <div className="flex flex-col mb-2">
                  <h3 className="font-semibold text-lg text-gray-800 dark:text-gray-200">Status: {task.status}</h3>
                  <span className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                    {formatDistanceToNow(new Date(task.created_at), { addSuffix: true })}
                  </span>
                </div>
                <p className="text-gray-700 dark:text-gray-300 mt-2">{task.details}</p>
              </div>
            ))
          ) : (
            <p className="text-center text-gray-500 dark:text-gray-400">No history available for this task.</p>
          )}
        </div>
      </SheetContent>
    </Sheet>
  );
}