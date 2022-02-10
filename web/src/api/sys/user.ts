import { defHttp } from '/@/utils/http/axios';

enum Api {
  GetPermCode = '/getPermCode',
}

export function getPermCode() {
  return defHttp.get<string[]>({ url: Api.GetPermCode });
}
