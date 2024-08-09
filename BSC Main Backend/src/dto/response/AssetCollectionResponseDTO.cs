namespace BSC_Main_Backend.dto.response;
using System.Collections.Generic;


/// <param name="Id">Id of AssetCollection</param>
/// <param name="Original">Id of original, non-baked collection</param>
/// <param name="Name">Name of collection, the original is named as "xxxx-original"</param>
/// <param name="Entries">Each GraphicAsset as part of this collection</param>
public record AssetCollectionResponseDTO(uint Id, uint Original, string Name, List<AssetCollectionEntryDTO> Entries);