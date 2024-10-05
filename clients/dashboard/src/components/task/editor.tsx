import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogClose,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import Editor from '@monaco-editor/react';

import { Button } from "@/components/ui/button"


import { Eye, Info, RefreshCw, RotateCw, Plus, X, type Volume1 } from 'lucide-react'; // Add this import for icons




/**
 * InputCode component to display task input in a dialog.
 * @param {Object} props - Component props
 * @param {any} props.input - Input data to display
 */
export function InputCode({ input }: { input: any }) {
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline" size="icon">
          <Eye className="h-4 w-4" />
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[700px] max-w-[95vw] w-full h-[60vh] p-0 overflow-hidden">
        <DialogHeader className="px-6 py-4 border-b">
          <div className="flex items-center justify-between">
            <DialogTitle className="text-xl font-semibold">Task Input</DialogTitle>
            <DialogClose asChild>
              <Button variant="ghost" size="icon" className="h-8 w-8">

              </Button>
            </DialogClose>
          </div>
        </DialogHeader>
        <div className="flex-grow h-full overflow-hidden">
          <Editor
            defaultLanguage="json"
            defaultValue={JSON.stringify(input, null, 2)}
            options={{
              readOnly: true,
              minimap: { enabled: false },
              scrollBeyondLastLine: false,
              fontSize: 14,
              lineNumbers: 'on',
              wordWrap: 'on',
              wrappingIndent: 'indent',
              automaticLayout: true,
            }}
            className="h-full"
          />
        </div>
      </DialogContent>
    </Dialog>
  );
}
