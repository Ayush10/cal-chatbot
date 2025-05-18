"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Calendar } from "@/components/ui/calendar"
import { CalendarIcon, X } from "lucide-react"
import { format } from "date-fns"
import { Dialog, DialogContent, DialogTrigger } from "@/components/ui/dialog"

interface SimpleDatePickerProps {
  onDateSelect: (date: Date | null) => void
  selectedDate: Date | null
}

export function SimpleDatePicker({ onDateSelect, selectedDate }: SimpleDatePickerProps) {
  const [open, setOpen] = useState(false)
  const [tempDate, setTempDate] = useState<Date | undefined>(selectedDate || undefined)

  const handleSelect = (date: Date | undefined) => {
    setTempDate(date)
  }

  const handleSave = () => {
    onDateSelect(tempDate || null)
    setOpen(false)
  }

  const handleClear = () => {
    onDateSelect(null)
  }

  return (
    <div className="relative">
      {selectedDate ? (
        <div className="flex items-center gap-2 p-2 bg-muted rounded-md">
          <CalendarIcon className="h-4 w-4" />
          <span className="text-xs">{format(selectedDate, "PPP")}</span>
          <Button variant="ghost" size="icon" className="h-5 w-5 rounded-full" onClick={handleClear}>
            <X className="h-3 w-3" />
          </Button>
        </div>
      ) : (
        <Dialog open={open} onOpenChange={setOpen}>
          <DialogTrigger asChild>
            <Button variant="ghost" size="icon" className="h-8 w-8">
              <CalendarIcon className="h-4 w-4" />
              <span className="sr-only">Open calendar</span>
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[425px]">
            <div className="p-3">
              <Calendar mode="single" selected={tempDate} onSelect={handleSelect} initialFocus />
              <div className="mt-4 flex justify-end gap-2">
                <Button variant="outline" onClick={() => setOpen(false)}>
                  Cancel
                </Button>
                <Button onClick={handleSave}>Select</Button>
              </div>
            </div>
          </DialogContent>
        </Dialog>
      )}
    </div>
  )
}
