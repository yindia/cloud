import { TaskManagementService } from "@buf/evalsocket_cloud.connectrpc_es/cloud/v1/cloud_connect"
import {
  createConnectTransport,
} from '@connectrpc/connect-web'
import {
  createPromiseClient,
} from '@connectrpc/connect'

// Update the TaskInterface to include a type field
export interface TaskInterface {
  id: number;
  name: string;
  type: string; // Add this line
  payload: {
    parameters: any;
  };
  retries: number;
  status: number;

  description: string
}


// Define the type based on the proto structure
export interface TaskHistory {
  id: number;
  details: string;
  status: string;
  created_at: string;
}


// Change this to a function that creates and returns the client
export function createTaskManagementClient() {
  return createPromiseClient(
    TaskManagementService,
    createConnectTransport({
      baseUrl: process.env.NEXT_PUBLIC_SERVER_ENDPOINT || "http://task:80"
    })
  );
}

// Add this function to convert status number to string
export function getStatusString(status: number): string {
  switch (status) {
    case 0:
      return "QUEUED";
    case 1:
      return "RUNNING";
    case 2:
      return "FAILED";
    case 3:
      return "SUCCEEDED";
    default:
      return "UNKNOWN";
  }
}



// Add this function to get status color
export const getStatusColor = (status: string): string => {
  switch (status) {
    case "QUEUED": return "text-yellow-500";
    case "RUNNING": return "text-blue-500";
    case "FAILED": return "text-red-500";
    case "SUCCEEDED": return "text-green-500";
    default: return "text-gray-500";
  }
};
