import { SignupForm } from "@/components/AuthForms/signup-form"
import { APP_NAME } from "@/config/constants"
import { Metadata } from "next"

export const metadata: Metadata = {
  title: 'Sign up for ' + APP_NAME,
}

export default function SignupPage() {
  return (
    <div className="bg-muted flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <a href="#" className="flex items-center gap-2 self-center font-medium">
          <div className="bg-primary text-primary-foreground flex size-6 items-center justify-center rounded-md">
            <img src={"/logo/icon-signee.png"} />
          </div>
          {APP_NAME}
        </a>
        <SignupForm />
      </div>
    </div>
  )
}
