using System.Reflection.Metadata;

namespace ASP.NETCoreWebAPI.Models
{
    public class LOD
    {
        public int Id { get; set; } // Primary Key
        public uint DetailLevel { get; set; } // Attribute - uint32, represents the detail level (0 for max detail, 1 for 1/4 detail, etc.)
        public required Blob Blob { get; set; } // Attribute - Required, binary data representing the image at this detail level
        public required GraphicalAsset GraphicalAsset { get; set; } // Foreign Key (GraphicalAsset) - Required, the associated graphical asset
    }
}
