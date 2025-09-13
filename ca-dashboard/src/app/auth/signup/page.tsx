'use client'
import { SignupPage } from "@/components/AuthForms/signup-form"
import { APP_NAME } from "@/config/constants"
import Image from "next/image"



export default function Signup() {
  return (
    <div className="bg-muted flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <a href="#" className="flex items-center gap-2 self-center font-medium">
          <div className="bg-primary text-primary-foreground flex size-6 items-center justify-center rounded-md">
            <Image width={200} height={200} alt="nothing" src="/logo/icon-signee.png" />
          </div>
          {APP_NAME}
        </a>
        <SignupPage />
      </div>
    </div>
  )
}
