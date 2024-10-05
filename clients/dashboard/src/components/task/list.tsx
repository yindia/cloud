'use client';

import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableFooter,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Button } from "@/components/ui/button"


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
import {
  createConnectTransport,
} from '@connectrpc/connect-web'
import {
  createPromiseClient,
} from '@connectrpc/connect'

import { formatDistanceToNow } from 'date-fns'
import { TaskManagementService } from "@buf/evalsocket_cloud.connectrpc_es/cloud/v1/cloud_connect"
import { Task, TaskType, CreateTaskRequest, Payload } from '@buf/evalsocket_cloud.bufbuild_es/cloud/v1/cloud_pb'; // Adjust this path if needed
import { useEffect, useState } from "react";  // Import useState
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Eye, Info, RefreshCw, RotateCw, Plus, X, type Volume1 } from 'lucide-react'; // Add this import for icons
import { CreateTask } from './create';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"

import { TaskHistory as TaskHistoryProto } from '@buf/evalsocket_cloud.bufbuild_es/cloud/v1/cloud_pb'; // Adjust this path if needed
import { TaskStatusEnum } from "@/cloud/v1/cloud_pb";
import { InputCode } from "@/components/task/editor";
import { createTaskManagementClient, getStatusColor, getStatusString, type TaskInterface } from "@/components/task/util";
import { TaskHistoryList } from "@/components/task/history";


/**
 * TaskList component to display a list of tasks.
 */
export function TaskList() {
  const [tasks, setTasks] = useState<TaskInterface[]>([]);  // State for tasks
  const [isOpen, setIsOpen] = useState(false); // Modal open state
  const [selectedTaskId, setSelectedTaskId] = useState<number>(0); // Selected task ID
  const [statusCounts, setStatusCounts] = useState({
    QUEUED: 0,
    RUNNING: 0,
    FAILED: 0,
    SUCCEEDED: 0,
    UNKNOWN: 0,
  });
  const [isCreateTaskOpen, setIsCreateTaskOpen] = useState(false);
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [statusFilter, setStatusFilter] = useState<string>("ALL");
  const [selectedTask, setSelectedTask] = useState<TaskInterface | null>(null);
  const [typeFilter, setTypeFilter] = useState<string>("ALL");
  const [currentPage, setCurrentPage] = useState(1);
  const [tasksPerPage, setTasksPerPage] = useState(10);
  const [hasMoreTasks, setHasMoreTasks] = useState(true);


  // Add this useEffect to fetch status counts every 5 seconds
  useEffect(() => {
    const fetchStatusCounts = async () => {
      try {
        const response = await createTaskManagementClient().getStatus({});
        const newStatusCounts: any = {};

        if ('statusCounts' in response) {
          Object.entries(response.statusCounts).forEach(([status, count]) => {
            newStatusCounts[getStatusString(parseInt(status))] = Number(count);
          });
        }
        setStatusCounts(newStatusCounts);
      } catch (error) {
        console.error('Error fetching status counts:', error);
        toast.error('Failed to fetch status counts. Please try again later.');
      }
    };

    fetchStatusCounts(); // Initial fetch
    const intervalId = setInterval(fetchStatusCounts, 5000); // Refresh every 5 seconds

    return () => clearInterval(intervalId);
  }, []);

  const fetchTasks = async (page: number) => {
    setIsRefreshing(true);
    try {
      const offset = (page - 1) * tasksPerPage;
      const request: any = {
        limit: tasksPerPage,
        offset: offset,
      };

      // Add status filter if it's not "ALL"
      if (statusFilter != "ALL") {
        request.status = TaskStatusEnum[statusFilter as keyof typeof TaskStatusEnum];
      } else {
        request.status = TaskStatusEnum.ALL
      }

      // Add type filter if it's not "ALL"
      if (typeFilter != "ALL") {
        request.type = typeFilter;
      }

      const response = await createTaskManagementClient().listTasks(request);

      const fetchedTasks = response.tasks.map((task: any) => ({
        id: task.id,
        name: task.name,
        type: task.type,
        description: task.description,
        payload: {
          parameters: task?.payload?.parameters || {},
        },
        status: task.status,
        retries: task.retries,
      }));

      setTasks(fetchedTasks);
      setHasMoreTasks(fetchedTasks.length === tasksPerPage);
    } catch (error) {
      console.error('Error fetching tasks:', error);
      toast.error('Failed to refresh tasks. Please try again later.');
    } finally {
      setIsRefreshing(false);
    }
  };

  // Update this useEffect to refetch tasks when filters change
  useEffect(() => {
    fetchTasks(1); // Reset to first page when filters change
    setCurrentPage(1);
  }, [statusFilter, typeFilter, tasksPerPage]);

  // Separate useEffect for periodic refresh
  useEffect(() => {
    const intervalId = setInterval(() => {
      fetchTasks(currentPage);
    }, 10000); // Fetch tasks every 10 seconds

    return () => clearInterval(intervalId);
  }, [currentPage, statusFilter, typeFilter, tasksPerPage]);

  const handleTaskClick = (taskId: number) => {
    setIsOpen(true); // Open task details modal
    setSelectedTaskId(taskId); // Set selected task ID
  };


  const handleCreateTask = async (newTask: TaskInterface) => {
    console.log(newTask, "===>");
  };

  const filteredTasks = tasks.filter((task) => {
    const statusMatch = statusFilter === "ALL" || getStatusString(task.status) === statusFilter;
    const typeMatch = typeFilter === "ALL" || task.type === typeFilter;
    return statusMatch && typeMatch;
  });


  // Get unique task types
  const taskTypes: any[] = ["ALL", "run_query", "send_email"];

  const handleNextPage = () => {
    if (hasMoreTasks) {
      setCurrentPage(prevPage => prevPage + 1);
    }
  };

  const handlePreviousPage = () => {
    if (currentPage > 1) {
      setCurrentPage(prevPage => prevPage - 1);
    }
  };

  useEffect(() => {
    fetchTasks(currentPage);
  }, [currentPage, statusFilter, typeFilter, tasksPerPage]);

  return (
    <>
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-2xl font-bold">Tasks Management</h2>
        <div className="space-x-2 flex items-center">
          <Select value={statusFilter} onValueChange={setStatusFilter}>
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Filter by status" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="ALL">All Statuses</SelectItem>
              <SelectItem value="QUEUED">Queued</SelectItem>
              <SelectItem value="RUNNING">Running</SelectItem>
              <SelectItem value="FAILED">Failed</SelectItem>
              <SelectItem value="SUCCEEDED">Succeeded</SelectItem>
              <SelectItem value="UNKNOWN">Unknown</SelectItem>
            </SelectContent>
          </Select>
          <Select value={typeFilter} onValueChange={setTypeFilter}>
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Filter by type" />
            </SelectTrigger>
            <SelectContent>
              {taskTypes.map((type) => (
                <SelectItem key={type} value={type}>
                  {type === "ALL" ? "All Types" : type}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <Button
            onClick={() => setIsCreateTaskOpen(true)}
            variant="outline"
            size="sm"
          >
            <Plus className="h-4 w-4 mr-2" />
            Add Task
          </Button>
          <Button
            onClick={() => fetchTasks(currentPage)}
            disabled={isRefreshing}
            variant="outline"
            size="sm"
          >
            <RotateCw className={`h-4 w-4 mr-2 ${isRefreshing ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
        {Object.entries(statusCounts).map(([status, count]) => (
          <Card key={status}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">{status}</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{count}</div>
            </CardContent>
          </Card>
        ))}
      </div>

      <Table>
        <TableCaption>A list of your recent tasks.</TableCaption>
        <TableHeader>
          <TableRow>
            <TableHead className="w-[100px]">ID</TableHead>
            <TableHead>Name</TableHead>
            <TableHead>Type</TableHead>
            <TableHead>Description</TableHead>
            <TableHead>Status</TableHead>
            <TableHead className="text-right">Action</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {filteredTasks.map((task: TaskInterface) => (
            <TableRow key={task.id}>
              <TableCell className="font-medium">{task.id}</TableCell>
              <TableCell>{task.name}</TableCell>
              <TableCell>{task.type}</TableCell>
              <TableCell>{task.description}</TableCell>
              <TableCell>
                <span className={`font-medium ${getStatusColor(getStatusString(task.status))}`}>
                  {getStatusString(task.status)}
                </span>
              </TableCell>
              <TableCell className="text-right">
                <div className="flex justify-end space-x-2">
                  <InputCode input={task?.payload?.parameters ?? ""} />
                  <Button
                    variant="outline"
                    size="icon"
                    onClick={() => handleTaskClick(task.id)}
                  >
                    <Info className="h-4 w-4" />
                  </Button>

                </div>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>

      <div className="flex justify-between items-center mt-4">
        <div>
          <span className="text-sm text-gray-700 dark:text-gray-400">
            Page {currentPage}
          </span>
        </div>
        <div className="space-x-2">
          <Button
            onClick={handlePreviousPage}
            disabled={currentPage === 1 || isRefreshing}
            variant="outline"
            size="sm"
          >
            Previous
          </Button>
          <Button
            onClick={handleNextPage}
            disabled={!hasMoreTasks || isRefreshing}
            variant="outline"
            size="sm"
          >
            Next
          </Button>
        </div>
      </div>

      {isOpen && (
        <TaskHistoryList
          task={tasks.find(t => t.id === selectedTaskId)!}
          isOpen={isOpen}
          setOpen={setIsOpen}
        />
      )}

      <CreateTask
        isOpen={isCreateTaskOpen}
        setOpen={setIsCreateTaskOpen}
        onCreateTask={handleCreateTask}
        initialTask={selectedTask}
      />
    </>
  );
}


