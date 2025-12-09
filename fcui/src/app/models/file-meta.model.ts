export interface FileMeta {
  id: number;
  filename: string;
  size: number;
  checksum: string;
  mime: string;
  path: string;
  created_at: string;
}

export interface FileMetaPage {
  items: FileMeta[];
  hasMore: boolean;
  nextCursor: string;
}