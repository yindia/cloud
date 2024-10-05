'use client';

import Head from 'next/head';
import * as React from 'react';
import '@/lib/env';
import { Toaster, toast } from 'react-hot-toast';

/**
 * SVGR Support
 * Caveat: No React Props Type.
 *
 * You can override the next-env if the type is important to you
 * @see https://stackoverflow.com/questions/68103844/how-to-override-next-js-svg-module-declaration
 */
import Logo from '~/svg/Logo.svg';
import clsx from 'clsx';
import { CreateTask } from '@/components/task/create';
import { TaskList } from '@/components/task/list';

// !STARTERCONF -> Select !STARTERCONF and CMD + SHIFT + F
// Before you begin editing, follow all comments with `STARTERCONF`,
// to customize the default configuration.

export default function HomePage() {
  const showToast = () => {
    toast.success('Hello, React Toast!');
  };

  const [mode, setMode] = React.useState<'dark' | 'light'>('light');

  return (
    <main className={clsx("flex flex-col items-center p-6 bg-gray-100 min-h-screen")}>
      <div className={clsx("bg-white shadow-md rounded-lg p-6 w-full max-w-4xl")}>
        <TaskList />
      </div>
    </main>
  );
}
