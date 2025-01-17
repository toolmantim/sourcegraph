import { useEffect, useMemo } from 'react'

import create from 'zustand'

import { Client, createClient, ClientInit } from '@sourcegraph/cody-shared/src/chat/client'
import { ChatMessage } from '@sourcegraph/cody-shared/src/chat/transcript/messages'
import { isErrorLike } from '@sourcegraph/common'

import { eventLogger } from '../tracking/eventLogger'

interface CodyChatStore {
    client: Client | null
    config: ClientInit['config'] | null
    messageInProgress: ChatMessage | null
    transcript: ChatMessage[]
    repo: string
    filePath: string
    setClient: (client: Client | null) => void
    setConfig: (config: ClientInit['config']) => void
    setMessageInProgress: (message: ChatMessage | null) => void
    setTranscript: (transcript: ChatMessage[]) => void
    initializeClient: (config: Required<ClientInit['config']>) => void
    onSubmit: (text: string) => void
    onReset: () => void
}

export const useChatStoreState = create<CodyChatStore>((set, get): CodyChatStore => {
    const onSubmit = (text: string): void => {
        const { client, repo, filePath } = get()
        if (client && !isErrorLike(client)) {
            eventLogger.log('web:codySidebar:submit', {
                repo,
                path: filePath,
                text,
            })
            client.submitMessage(text)
        }
    }

    const onReset = (): void => {
        const { initializeClient, config } = get()
        if (config) {
            initializeClient(config as Required<ClientInit['config']>)
        }
    }

    return {
        client: null,
        messageInProgress: null,
        config: null,
        transcript: [],
        filePath: '',
        repo: '',
        setClient: client => set({ client }),
        setConfig: config => set({ config }),
        setMessageInProgress: message => set({ messageInProgress: message }),
        setTranscript: transcript => set({ transcript }),
        initializeClient: (config: Required<ClientInit['config']>): void => {
            set({ messageInProgress: null, transcript: [], repo: config.codebase, config })
            createClient({
                config,
                setMessageInProgress: message => set({ messageInProgress: message }),
                setTranscript: transcript => set({ transcript }),
            })
                .then(client => {
                    set({ client })
                })
                .catch(error => {
                    eventLogger.log('web:codySidebar:clientError', { repo: config?.codebase })
                    set({ client: error })
                })
        },

        onSubmit,
        onReset,
    }
})

export const useChatStore = (isCodyEnabled: boolean, repoName: string): CodyChatStore => {
    const store = useChatStoreState()

    const config = useMemo<Required<ClientInit['config']>>(
        () => ({
            serverEndpoint: window.location.origin,
            useContext: 'embeddings',
            codebase: repoName,
            accessToken: null,
        }),
        [repoName]
    )

    const { initializeClient } = store
    useEffect(() => {
        if (!isCodyEnabled) {
            return
        }

        initializeClient(config)
    }, [config, initializeClient, isCodyEnabled])

    return store
}
