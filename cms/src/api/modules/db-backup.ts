import request from '../request'
import type {
    GetDbBackupCfgRes,
    ListDbBackupFilesRes,
    RestoreDbBackupReq,
    RestoreDbBackupRes,
    RunDbBackupNowRes,
    SaveDbBackupCfgReq,
    SaveDbBackupCfgRes,
} from '@/types/api'

export const dbBackupApi = {
    getCfg: () => request.post<GetDbBackupCfgRes>('/dbBackup/getDbBackupCfg', {}),
    saveCfg: (data: SaveDbBackupCfgReq) => request.post<SaveDbBackupCfgRes>('/dbBackup/saveDbBackupCfg', data),
    runNow: () => request.post<RunDbBackupNowRes>('/dbBackup/runDbBackupNow', {}),
    listFiles: () => request.post<ListDbBackupFilesRes>('/dbBackup/listDbBackupFiles', {}),
    restore: (data: RestoreDbBackupReq) => request.post<RestoreDbBackupRes>('/dbBackup/restoreDbBackup', data),
}

export default dbBackupApi
