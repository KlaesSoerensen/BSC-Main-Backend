using System.Reflection.Metadata;

namespace ASP.NETCoreWebAPI.Models
{
    public class GraphicalAsset
    {
        public int Id { get; set; } // Primary Key
        public required string Type { get; set; } // Attribute - Required, represents the type of the asset (e.g., "png", "jpg", "webp", "svg")
        public required string UseCase { get; set; } // Attribute - Required, describes the use case (e.g., "icon", "environment", "player")
        public bool HasLODs { get; set; } // Attribute - Indicates if the asset has Levels of Detail (LODs)
        public Blob? Blob { get; set; } // Attribute - Nullable, represents the binary data of the asset, immediately available if no LODs
        public int Width { get; set; } // Attribute - Width of the asset
        public int Height { get; set; } // Attribute - Height of the asset
    }
}
